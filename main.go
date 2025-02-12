package main

import (
	"context"
	"log"
	"os"

	"github.com/ghofaralhasyim/be-appointment-system/internal/config"
	"github.com/ghofaralhasyim/be-appointment-system/internal/middleware"
	"github.com/ghofaralhasyim/be-appointment-system/internal/routes"
	"github.com/ghofaralhasyim/be-appointment-system/pkg/database"
	"github.com/ghofaralhasyim/be-appointment-system/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	echo := echo.New()

	v := validator.New()
	echo.Validator = &CustomValidator{validator: v}

	echo.Use(middleware.CORSMiddleware)

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.InitDbConnection()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	redisConf, err := config.NewRedisConfig()
	if err != nil {
		log.Fatalf("Could not connect to redis server: %v", err)
	}
	redisClient := database.NewRedisClient(redisConf)
	_, err = redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("redis connection failed: %v", err)
	}
	log.Println("check: redis connected")

	routes.SetupRoutes(echo, db, redisClient)

	if os.Getenv("STAGE_STATUS") == "production" {
		utils.StartServerWithGracefulShutdown(echo)
	} else {
		utils.StartServer(echo)
	}
}
