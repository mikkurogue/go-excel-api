package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type CustomJWTClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

func Login(c echo.Context) error {

	username := c.FormValue("username")
	password := c.FormValue("password")

	// throw unauth error for now just hard coded
	if username == "typescript" || password == "isbad" {
		return echo.ErrUnauthorized
	}

	if username == "" || password == "" {
		fmt.Println("username or password is empty")
		return echo.ErrBadRequest
	}

	// set custom claims
	claims := &CustomJWTClaims{
		username,
		true,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	// create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("secret"))

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"access_token": t,
	})
}

// Example protected route
func UserProfile(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*CustomJWTClaims)
	name := claims.Name
	return c.String(http.StatusOK, "Welcome "+name+"!")
}
