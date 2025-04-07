package main

import (
	"expvar"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/Althaf66/Appointr/internal/auth"
	"github.com/Althaf66/Appointr/internal/db"
	// "github.com/Althaf66/Appointr/internal/env"
	"github.com/Althaf66/Appointr/internal/mailer"
	"github.com/Althaf66/Appointr/internal/store"
	"github.com/Althaf66/Appointr/internal/websocket"
	"github.com/joho/godotenv"
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
	godotenv.Load()
	cfg := config{
		addr:        os.Getenv("ADDR"),
		apiUrl:      os.Getenv("EXTERNAL_URL"),
		frontendURL: os.Getenv("FRONTEND_URL"),
		env:         os.Getenv("ENV"),
		db: dbConfig{
			addr:         os.Getenv("DB_ADDR"),
			maxOpenConns: 30,
			maxIdleConns: 30,
			maxIdleTime:  os.Getenv("DB_MAX_IDLE_TIME"),
		},
		mail: mailconfig{
			exp:       time.Hour * 24 * 3,
			fromEmail: os.Getenv("FROM_EMAIL"),
			mailTrap: mailTrapConfig{
				apiKey: os.Getenv("MAILTRAP_API_KEY"),
			},
		},

		auth: authConfig{
			basic: basicConfig{
				username: os.Getenv("AUTH_BASIC_USER"),
				password: os.Getenv("AUTH_BASIC_PASS"),
			},
			token: tokenConfig{
				secret: os.Getenv("AUTH_TOKEN_SECRET"),
				exp:    time.Hour * 24 * 3,
				iss:    "appointr",
			},
		},
		stripeKey:     os.Getenv("STRIPE_KEY"),
		stripeWebhook: os.Getenv("STRIPE_WEBHOOK"),
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
