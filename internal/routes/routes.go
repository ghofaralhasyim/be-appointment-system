package routes

import (
	"database/sql"

	"github.com/ghofaralhasyim/be-appointment-system/internal/http"
	"github.com/ghofaralhasyim/be-appointment-system/internal/middleware"
	"github.com/ghofaralhasyim/be-appointment-system/internal/repositories"
	"github.com/ghofaralhasyim/be-appointment-system/internal/services"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, db *sql.DB, redisClient *redis.Client) {
	apiV1 := e.Group("/v1")

	redisRepo := repositories.NewRedisRepository(redisClient)

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo, redisRepo)
	userHandler := http.NewUserHandler(userService)
	apiV1.POST("/auth/login", userHandler.Login)
	apiV1.POST("/auth/refresh", userHandler.RefreshToken, middleware.AuthMiddleware(redisRepo))

	appointmentRepo := repositories.NewAppointmentRepository(db)
	appointmentService := services.NewAppointmentService(appointmentRepo)
	appointmentHandler := http.NewAppointmentHandler(appointmentService)
	apiV1.POST("/appointment", appointmentHandler.CreateAppointment, middleware.AuthMiddleware(redisRepo))
}
