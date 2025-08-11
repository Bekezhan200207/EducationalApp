package handlers

import (
	"go-EdTech/docs"
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

	r.GET("/users/:uuid", usersHandlers.FindOne)
	r.POST("/users", usersHandlers.Create)
	r.GET("/users", usersHandlers.FindAll)
	r.PUT("/users/:uuid/:uuid", usersHandlers.Update)
	r.DELETE("/users/:uuid", usersHandlers.Delete)
	r.PATCH("/users/:uuid/:uuid/changePassword", usersHandlers.ChangePassword)
	r.PATCH("/users/:uuid/deactivate", usersHandlers.Deactivate)
	r.PATCH("/users/:uuid/activate", usersHandlers.Activate)

	r.GET("/lessons/:id", lessonsHandlers.FindById)
	r.GET("/lessons", lessonsHandlers.FindAll)
	r.POST("/lessons", lessonsHandlers.Create)
	r.PUT("lessons/:id", lessonsHandlers.Update)
	r.DELETE("lessons/:id", lessonsHandlers.Delete)

	r.GET("/subjects/:id", subjectsHandlers.FindById)
	r.GET("/subjects", subjectsHandlers.FindAll)
	r.POST("/subjects", subjectsHandlers.Create)
	r.PUT("/subjects/:id", subjectsHandlers.Update)
	r.DELETE("/subjects/:id", subjectsHandlers.Delete)

	r.GET("/courses/:id", CoursesHandlers.FindById)
	r.GET("/courses", CoursesHandlers.FindAll)
	r.POST("/courses", CoursesHandlers.Create)
	r.PUT("courses/:id", CoursesHandlers.Update)
	r.DELETE("courses/:id", CoursesHandlers.Delete)

	// Swagger
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", swagger.WrapHandler(swaggerfiles.Handler))
}
