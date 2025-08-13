package handlers

import (
	"go-EdTech/docs"
	"go-EdTech/middlewares"
	"go-EdTech/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerfiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(r *gin.Engine, conn *pgxpool.Pool) {

	usersRepository := repositories.NewUsersRepository(conn)
	lessonsRepository := repositories.NewLessonsRepository(conn)
	subjectsRepository := repositories.NewSubjectsRepository(conn)
	coursesRepository := repositories.NewCoursesRepository(conn)

	usersHandlers := NewUsersHandlers(usersRepository)
	lessonsHandlers := NewLessonsHandler(lessonsRepository)
	subjectsHandlers := NewSubjectsHandlers(subjectsRepository)
	CoursesHandlers := NewCoursesHandler(coursesRepository)
	authHandlers := NewAuthHandlers(usersRepository)

	authorized := r.Group("")
	authorized.Use(middlewares.AuthMiddleware)

	authorized.GET("/users/:uuid", usersHandlers.FindById)
	authorized.GET("/users", usersHandlers.FindAll)
	authorized.PUT("/users/:uuid/:uuid", usersHandlers.Update)
	authorized.DELETE("/users/:uuid", usersHandlers.Delete)
	authorized.PATCH("/users/:uuid/:uuid/changePassword", usersHandlers.ChangePassword)
	authorized.PATCH("/users/:uuid/deactivate", usersHandlers.Deactivate)
	authorized.PATCH("/users/:uuid/activate", usersHandlers.Activate)
	authorized.POST("/auth/signOut", authHandlers.SignOut)
	authorized.GET("/auth/userInfo", authHandlers.GetUserInfo)

	authorized.GET("/lessons/:id", lessonsHandlers.FindById)
	authorized.GET("/lessons", lessonsHandlers.FindAll)
	authorized.POST("/lessons", lessonsHandlers.Create)
	authorized.PUT("lessons/:id", lessonsHandlers.Update)
	authorized.DELETE("lessons/:id", lessonsHandlers.Delete)

	authorized.GET("/subjects/:id", subjectsHandlers.FindById)
	authorized.GET("/subjects", subjectsHandlers.FindAll)
	authorized.POST("/subjects", subjectsHandlers.Create)
	authorized.PUT("/subjects/:id", subjectsHandlers.Update)
	authorized.DELETE("/subjects/:id", subjectsHandlers.Delete)

	authorized.GET("/courses/:id", CoursesHandlers.FindById)
	authorized.GET("/courses", CoursesHandlers.FindAll)
	authorized.POST("/courses", CoursesHandlers.Create)
	authorized.PUT("courses/:id", CoursesHandlers.Update)
	authorized.DELETE("courses/:id", CoursesHandlers.Delete)

	unauthorized := r.Group("")
	unauthorized.POST("/auth/signIn", authHandlers.SignIn)
	unauthorized.POST("/users", usersHandlers.RegisterNewUser)

	// Swagger
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", swagger.WrapHandler(swaggerfiles.Handler))
}
