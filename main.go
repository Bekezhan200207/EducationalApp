package main

import (
	"context"
	"fmt"
	"go-EdTech/config"
	"go-EdTech/handlers"
	"go-EdTech/logger"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// @title 			EdTech API
// @version 		1.0
// @description		this is a Educational Application project
// @termsOfService 	http://swagger.io/terms/
//
// @contact.name 	API Support
// @contact.url	 	http://www.swagger.io/support
// @contact.email 	support@swagger.io
//
// @license.name 	Apache 2.0
// @license.url 	http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host 			https://ilessons.cloud/go/api
// @BasePath 		/
//
// @externalDocs.description 	OpenAPI
// @externalDocs.url 			https://swagger.io/resources/open.api/

func main() {
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)

	logger := logger.GetLogger()

	defer func() {
		if r := recover(); r != nil {
			logger.Error("Application crashed!", zap.Any("error", r))
		}
	}()

	r.Use(
		ginzap.Ginzap(logger, time.RFC3339, true),
		ginzap.RecoveryWithZap(logger, true),
	)

	corsConfig := cors.Config{
		AllowOrigins:     []string{"ilessons.cloud/go/api"},
		AllowHeaders:     []string{"Content-Type, Authorization"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	}

	r.Use(cors.New(corsConfig))

	logger.Info("Loading configuration...")
	err := loadConfig()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	logger.Info("Connecting to database...")
	conn, err := connectToDb()
	if err != nil {
		logger.Fatal("Database connection failed", zap.Error(err))
	}

	r.Use(func(c *gin.Context) {
		c.Set("db", conn)
		c.Next()
	})

	handlers.SetupRoutes(r, conn)

	logger.Info("Application starting...", zap.String("host", config.Config.AppHost))
	if err := r.Run(config.Config.AppHost); err != nil {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}

func loadConfig() error {
	// Указываем путь к .env файлу
	viper.SetConfigFile(".env")

	// Загружаем переменные из .env, если он есть (необязательно)
	_ = viper.ReadInConfig() // не падаем, если файла нет

	// Читаем переменные окружения (например, из Railway)
	viper.AutomaticEnv()
	// Bind specific environment variables
	viper.BindEnv("APP_HOST")
	viper.BindEnv("DB_CONNECTION_STRING")
	viper.BindEnv("JWT_SECRET_KEY")
	viper.BindEnv("JWT_EXPIRE_DURATION")

	// Мапим переменные в структуру
	var mapConfig config.MapConfig
	err := viper.Unmarshal(&mapConfig)
	if err != nil {
		return err
	}

	// Validate required configuration
	if mapConfig.DbConnectionString == "" {
		return fmt.Errorf("DB_CONNECTION_STRING environment variable is required")
	}
	if mapConfig.JwtSecretKey == "" {
		return fmt.Errorf("JWT_SECRET_KEY environment variable is required")
	}
	if mapConfig.AppHost == "" {
		mapConfig.AppHost = ":8081" // Default value
	}

	config.Config = &mapConfig
	return nil

}

func connectToDb() (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(context.Background(), config.Config.DbConnectionString)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return conn, nil
}
