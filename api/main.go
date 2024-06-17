package main

import (
	"go-backend/handlers"
	"log"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

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

	r := e.Group("/auth")
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(handlers.CustomJWTClaims)
		},
		SigningKey: []byte("secret"),
	}

	r.Use(echojwt.WithConfig(config))
	r.GET("/user/profile", handlers.UserProfile)

	r.POST("/upload/excel", handlers.UploadExcel)

	r.GET("/process/all", handlers.GetAllProcesses)
	r.GET("/process/:id", handlers.GetProcessById)

	e.Logger.Fatal(e.Start(":1323"))
}
