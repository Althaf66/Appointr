package main

import (
	"expvar"
	"log"
	"runtime"
	"time"

	"github.com/Althaf66/Appointr/internal/auth"
	"github.com/Althaf66/Appointr/internal/db"
	"github.com/Althaf66/Appointr/internal/env"
	"github.com/Althaf66/Appointr/internal/mailer"
	"github.com/Althaf66/Appointr/internal/store"
	"github.com/Althaf66/Appointr/internal/websocket"
	"go.uber.org/zap"
)

const version = "1.0.0"

//	@title			Appointr API
//	@version		1.0
//	@description	This is a Appointr server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host						localhost:8080
// @BasePath					/v1
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	cfg := config{
		addr:        env.GetString("ADDR", ":8080"),
		apiUrl:      env.GetString("EXTERNAL_URL", "localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:5173"),
		env:         env.GetString("ENV", "development"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/appointr?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		mail: mailconfig{
			exp:       time.Hour * 24 * 3,
			fromEmail: env.GetString("FROM_EMAIL", "althafas@althafasharaf.in"),
			mailTrap: mailTrapConfig{
				apiKey: env.GetString("MAILTRAP_API_KEY", "b1baf93cac0045d35b7f2299c2398d1a"),
			},
		},

		auth: authConfig{
			basic: basicConfig{
				username: env.GetString("AUTH_BASIC_USER", "admin"),
				password: env.GetString("AUTH_BASIC_PASS", "admin"),
			},
			token: tokenConfig{
				secret: env.GetString("AUTH_TOKEN_SECRET", "unknown"),
				exp:    time.Hour * 24 * 3,
				iss:    "appointr",
			},
		},
		stripeKey:     env.GetString("STRIPE_KEY", "sk_test_51NpubpSF2Iapc3CYMensOJgnQjo4anfwi9MNLOFIjNkOBYRxzEP8gMctadHwISPfAERy31iKNejTs50cRCu1bCxV00NycUfZ06"),
		stripeWebhook: env.GetString("STRIPE_WEBHOOK", "whsec_5048a44b5392d726f42c8c6ed46045c2f062db8d551ed0dd1ff239c4bcbdf26d"),
	}

	// logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	//database
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	logger.Info("database connection pool established")

	store := store.NewPostgresStorage(db)
	wsManager := websocket.NewWebSocketManager(store)

	mailtrap, err := mailer.NewMailTrapClient(cfg.mail.mailTrap.apiKey, cfg.mail.fromEmail)
	if err != nil {
		logger.Fatal(err)
	}
	jwtAuthenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.iss, cfg.auth.token.iss)

	app := application{
		config:        cfg,
		store:         store,
		logger:        logger,
		mailer:        mailtrap,
		authenticator: jwtAuthenticator,
		wsManager:     wsManager,
	}

	expvar.NewString("version").Set(version)
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	mux := app.mount()
	log.Fatal(app.run(mux))

}
