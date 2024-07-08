package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go-backend/db"
	"go-backend/handlers"
	"log"
	"os"
)

func main() {
	e := echo.New()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading environment variables")
	}

	database, err := db.Init(os.Getenv("TURSO_DB_NAME"), os.Getenv("TURSO_AUTH_TOKEN"))
	if err != nil {
		log.Fatal(err.Error())
	}

	activeDatabase := database.CreateConnection()

	usrs := db.QueryUsers(activeDatabase)

	fmt.Println(usrs)

	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Error code 0x901: Request time out - Contact dev team if this happens frequently",
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			log.Println(c.Path())
		},
		Timeout: 0,
	}))

	e.Use(middleware.BodyLimit("70M"))

	e.GET("/core", handlers.Core)

	e.POST("/user/login", handlers.Login)

	// for now un-authenticate these routes for easier testing.
	e.POST("/upload/excel", handlers.UploadExcel)
	e.GET("/process/all", handlers.GetAllProcesses)
	e.GET("/process/:id", handlers.GetProcessById)

	e.DELETE("/process/:id", handlers.DeleteProcess)

	r := e.Group("/auth")
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(handlers.CustomJWTClaims)
		},
		SigningKey: []byte("secret"),
	}

	r.Use(echojwt.WithConfig(config))
	r.GET("/user/profile", handlers.UserProfile)

	e.Logger.Fatal(e.Start(":1323"))
}
