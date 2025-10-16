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
	roleRepository := repositories.NewRoleRepository(conn)
	sessionsRepository := repositories.NewSessionsRepository(conn)

	usersHandlers := NewUsersHandlers(usersRepository)
	lessonsHandlers := NewLessonsHandler(lessonsRepository)
	subjectsHandlers := NewSubjectsHandlers(subjectsRepository)
	CoursesHandlers := NewCoursesHandler(coursesRepository)
	authHandlers := NewAuthHandler(usersRepository, sessionsRepository, roleRepository)

	unauthorized := r.Group("")


	authorized := r.Group("/")
	authorized.Use(middlewares.AuthMiddleware(sessionsRepository, usersRepository, roleRepository))

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/signup", authHandlers.SignUp)
		authGroup.POST("/login", authHandlers.Login)
		authGroup.POST("/logout", authHandlers.Logout)
		authGroup.POST("/refresh", authHandlers.Refresh)
	}	

	authorized.GET("/users/:uuid", usersHandlers.FindById)
	authorized.GET("/users", usersHandlers.FindAll)
	authorized.PUT("/users/:uuid", usersHandlers.Update)
	authorized.DELETE("/users/:uuid", usersHandlers.Delete)
	unauthorized.PATCH("/users/:uuid/changePassword", usersHandlers.ChangePassword)  //New logic needed. Restoration by email
	authorized.PATCH("/users/:uuid/deactivate", usersHandlers.Deactivate) //no functionality 
	authorized.PATCH("/users/:uuid/activate", usersHandlers.Activate) //no functionality

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

	// Swagger
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", swagger.WrapHandler(swaggerfiles.Handler))
}
