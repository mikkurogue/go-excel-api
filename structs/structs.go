package structs

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

func kys() {
	fmt.Print("kys")
}
