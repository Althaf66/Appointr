package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Althaf66/Appointr/docs"
	"github.com/Althaf66/Appointr/internal/auth"
	"github.com/Althaf66/Appointr/internal/env"
	"github.com/Althaf66/Appointr/internal/mailer"
	"github.com/Althaf66/Appointr/internal/store"
	"github.com/Althaf66/Appointr/internal/websocket"
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type application struct {
	config        config
	store         store.Storage
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Authenticator
	wsManager     *websocket.WebSocketManager
}

type config struct {
	addr          string
	env           string
	frontendURL   string
	apiUrl        string
	db            dbConfig
	mail          mailconfig
	auth          authConfig
	stripeKey     string
	stripeWebhook string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type authConfig struct {
	basic basicConfig
	token tokenConfig
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

type basicConfig struct {
	username string
	password string
}

type mailconfig struct {
	exp       time.Duration
	mailTrap  mailTrapConfig
	fromEmail string
}

type mailTrapConfig struct {
	apiKey string
}

func (app *application) mount() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{env.GetString("CORS_ALLOWED_ORIGIN", "http://localhost:5173")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/ws", func(r chi.Router) {
		r.Get("/messages/{conversationID}", app.HandleWebSocket(app.wsManager))
	})

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthcheckHandler)
		r.With(app.BasicAuthMiddleware()).Get("/debug/vars", expvar.Handler().ServeHTTP)

		docsUrl := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(docsUrl),
		))
		r.Route("/expertise", func(r chi.Router) {
			r.Get("/", app.getExpertiseHandler)
			r.Post("/create", app.createExpertiseHandler)
			r.Get("/{expertiseID}", app.getExpertiseHandlerByID)
			r.Patch("/{expertiseID}", app.updateExpertiseHandler)
			r.Delete("/{expertiseID}", app.deleteExpertiseHandler)
		})
		r.Route("/discipline", func(r chi.Router) {
			r.Get("/", app.getDisciplineHandler)
			r.Post("/create", app.createDisciplineHandler)
			r.Get("/{disciplineField}", app.getDisciplineHandlerByField)
			r.Group(func(r chi.Router) {
				r.Patch("/{disciplineID}", app.updateDisciplineHandler)
				r.Delete("/{disciplineID}", app.deleteDisciplineHandler)
			})
		})
		r.Route("/users", func(r chi.Router) {
			// r.Get("/auth/google/login", app.googleLogin)
			// r.Get("/auth/google/callback", app.googleCallback)

			r.Put("/activate/{token}", app.activateUserHandler)

			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/", app.getUserHandler)
			})
		})
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
			r.Post("/token", app.createTokenHandler)
		})
		r.Post("/webhook", app.handleWebhook)
		r.Route("/payment", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/create-checkout-session", app.createCheckoutSession)
		})
		r.Route("/messages", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/conversations", app.createConversationHandler)
			r.Post("/conversations/{id}/messages", app.createMessageHandler)
			r.Get("/conversations", app.getConversationsHandler)
			r.Get("/conversations/{conversationID}", app.getConversationHandler)
			r.Get("/{conversationID}", app.getMessagesHandler)
			r.Put("/{conversationID}/read", app.markConversationReadHandler)
			r.Get("/unread", app.getUnreadCountHandler)
		})
		r.Route("/mentors", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/create", app.createMentorHandler)
			r.Get("/", app.getMentorsHandler)
			r.Get("/name/{mentorName}", app.getMentorByNameHandler)
			r.Get("/exp/{expertise}", app.getMentorByExpertiseHandler)
			r.Get("/dis/{discipline}", app.getMentorByDisciplineHandler)
			r.Get("/u/{userID}", app.getMentorByUserIDHandler)
			r.Group(func(r chi.Router) {
				r.Use(app.mentorContextMiddleware)
				r.Get("/{mentorID}", app.getMentorByIDHandler)
				r.Patch("/{mentorID}", app.updateMentorHandler)
				r.Delete("/{mentorID}", app.deleteMentorHandler)
			})
		})
		r.Route("/gigs", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/create", app.createGigHandler)
			r.Get("/", app.getAllGigsHandler)
			r.Get("/expertise/{expertise}", app.getGigsByExpertiseHandler)
			r.Get("/{gigID}", app.getGigByIDHandler)
			r.Patch("/{gigID}", app.updateGigHandler)
			r.Delete("/{gigID}", app.deleteGigHandler)
			// r.Group(func(r chi.Router) {
			// 	r.Use(app.gigContextMiddleware)
			// })
		})
		r.Route("/education", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/create", app.createEducationHandler)
			r.Get("/u/{id}", app.getEducationByUserIDHandler)
			r.Get("/{educationID}", app.getEducationByIDHandler)
			r.Patch("/{educationID}", app.updateEducationHandler)
			r.Delete("/{educationID}", app.deleteEducationHandler)
		})
		r.Route("/experience", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/create", app.createExperienceHandler)
			r.Get("/u/{id}", app.getExperienceByUserIDHandler)
			r.Get("/{experienceID}", app.getExperienceByIDHandler)
			r.Patch("/{experienceID}", app.updateExperienceHandler)
			r.Delete("/{experienceID}", app.deleteExperienceHandler)
		})
		r.Route("/socialmedia", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/create", app.createSocialMediaHandler)
			r.Get("/u/{id}", app.getSocialMediaByUserIDHandler)
			r.Get("/{socialMediaID}", app.getSocialMediaByIDHandler)
			r.Patch("/{socialMediaID}", app.updateSocialMediaHandler)
			r.Delete("/{socialMediaID}", app.deleteSocialMediaHandler)
		})
		r.Route("/workingat", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/create", app.createWorkingAtHandler)
			r.Get("/u/{id}", app.getWorkingAtByUserIDHandler)
			r.Get("/{workingatID}", app.getWorkingAtByIDHandler)
			r.Patch("/{workingatID}", app.updateWorkingAtHandler)
			r.Delete("/{workingatID}", app.deleteWorkingAtHandler)
		})
		r.Put("/meetings/paid/{meetingID}", app.updateMeetingPaidHandler)
		r.Route("/meetings", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/create", app.createMeetingHandler)
			r.Get("/", app.getAllMeetingsHandler)
			r.Get("/u/{userID}", app.getMeetingByUserIDHandler)
			r.Get("/mentor-not-confirm/{mentorID}", app.getMeetingMentorNotConfirmHandler)
			r.Get("/user-not-paid/{userID}", app.getMeetingUserNotPaidHandler)
			r.Get("/user-not-completed/{userID}", app.getMeetingUserNotCompletedHandler)
			r.Get("/mentor-not-completed/{mentorID}", app.getMeetingMentorNotCompletedHandler)
			r.Put("/confirm/{meetingID}", app.updateMeetingConfirmHandler)
			r.Put("/completed/{meetingID}", app.updateMeetingCompletedHandler)
			r.Delete("/{meetingID}", app.deleteMeetingHandler)
		})
		r.Route("/bookingslots", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/create", app.createBookingSlotHandler) // POST /bookingslots/create
			// r.Get("/", getBookingSlots)          // GET /bookingslots
			// r.Patch("/{id}", updateBookingSlot)  // PATCH /bookingslots/{id}
			// r.Delete("/{id}", deleteBookingSlot) // DELETE /bookingslots/{id}
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {

	docs.SwaggerInfo.Version = "0.0.1"
	docs.SwaggerInfo.Host = app.config.apiUrl
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.logger.Infow("signal caught", "signal", s.String())
		shutdown <- srv.Shutdown(ctx)
	}()

	app.logger.Info("server has started ", "addr", app.config.addr, " env:", app.config.env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	app.logger.Infow("server has stopped", "addr", app.config.addr, "env", app.config.env)

	return nil
}
