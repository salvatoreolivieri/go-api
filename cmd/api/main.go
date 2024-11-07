package main

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/salvatoreolivieri/go-api/internal/auth"
	"github.com/salvatoreolivieri/go-api/internal/db"
	"github.com/salvatoreolivieri/go-api/internal/env"
	"github.com/salvatoreolivieri/go-api/internal/mailer"
	"github.com/salvatoreolivieri/go-api/internal/store"
	"github.com/salvatoreolivieri/go-api/internal/store/cache"
	"go.uber.org/zap"
)

const version = "0.0.1"

//	@title			GopherSocial API
//	@description	API for GopherSocial, a social network for gohpers
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {

	config := config{
		addr:        env.GetString("ADDR", ":8000"),
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:8000"),
		frontendURL: env.GetString("FRONTEND_URL", "localhost:5173"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/socialnetwork?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		redisConfig: redisConfig{
			addr:     env.GetString("REDIS_ADDR", "localhost:6379"),
			password: env.GetString("REDIS_PASSWORD", ""),
			db:       env.GetInt("REDIS_DB", 0),
			enabled:  env.GetBool("REDIS_ENABLED", true),
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
			fromEmail:  env.GetString("FROM_EMAIL", "salvatore.olivieri07@gmail.com"),
			expiration: time.Hour * 24 * 3, // 3 days
		},
		auth: authConfig{
			basic: basicConfig{
				user: "admin",
				pass: "admin",
			},
			token: tokenConfig{
				secret:     env.GetString("AUTH_TOKEN_SECRET", "example"),
				expiration: time.Hour * 24 * 3, // 3 days
				issuer:     "gophersocial",
			},
		},
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Database
	db, err := db.New(
		config.db.addr,
		config.db.maxOpenConns,
		config.db.maxIdleConns,
		config.db.maxIdleTime,
	)

	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("database connection pool established")

	// Cache
	var redisDb *redis.Client
	if config.redisConfig.enabled {
		redisDb = cache.NewRedisCLient(config.redisConfig.addr, config.redisConfig.password, config.redisConfig.db)
		logger.Info("redis cache connection established")
	}

	store := store.NewStorage(db)
	cacheStorage := cache.NewRedisStorage(redisDb)

	mailer := mailer.NewSendgrid(config.mail.sendGrid.apiKey, config.mail.fromEmail)

	authenticator := auth.NewJWTAuthenticator(
		config.auth.token.secret,
		config.auth.token.issuer,
		config.auth.token.issuer,
	)

	app := &application{
		config,
		store,
		cacheStorage,
		logger,
		mailer,
		authenticator,
	}

	// instantiate the handler
	mux := app.mount()

	logger.Fatal(app.run(mux))

}
