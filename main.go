package main

import (
	"context"
	"go-EdTech/config"
	"go-EdTech/docs"
	"go-EdTech/handlers"
	"go-EdTech/logger"
	"go-EdTech/repositories"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
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
// @host 			api.ozinshe.com
// @BasePath 		/
//
// @externalDocs.description 	OpenAPI
// @externalDocs.url 			https://swagger.io/resources/open.api/

func main() {
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)

	logger := logger.GetLogger()
	r.Use(
		ginzap.Ginzap(logger, time.RFC3339, true),
		ginzap.RecoveryWithZap(logger, true),
	)

	corsConfig := cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"*"},
		AllowMethods:    []string{"*"},
	}

	r.Use(cors.New(corsConfig))

	err := loadConfig()
	if err != nil {
		panic(err)
	}

	conn, err := connectToDb()
	if err != nil {
		panic("Could not connect to db")
	}

	usersRepository := repositories.NewUsersRepository(conn)
	usersHandlers := handlers.NewUsersHandlers(usersRepository)

	r.GET("/core/V1/user/profile/:uuid", usersHandlers.FindOne)                         //FWP
	r.POST("/core/V1/user/profile", usersHandlers.Create)                               //FWP
	r.GET("/core/V1/user/profile", usersHandlers.FindAll)                               //FWP
	r.PUT("/core/V1/user/profile/:uuid", usersHandlers.Update)                          //FWP
	r.DELETE("/core/V1/user/profile/:uuid", usersHandlers.Delete)                       //FWP
	r.PATCH("/core/V1/user/profile/:uuid/changePassword", usersHandlers.ChangePassword) //FWP
	r.PATCH("/core/V1/user/profile/:uuid/deactivate", usersHandlers.Deactivate)         //FWP
	r.PATCH("/core/V1/user/profile/:uuid/activate", usersHandlers.Activate)             //FWP

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", swagger.WrapHandler(swaggerfiles.Handler))

	logger.Info("Application starting")

	r.Run(config.Config.AppHost)
}

func loadConfig() error {
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	var mapConfig config.MapConfig
	err = viper.Unmarshal(&mapConfig)
	if err != nil {
		return err
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
