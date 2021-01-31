package jwtconfig

import (
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	models "linkaja.com/e-wallet/lib/base_models"
)

// CreateToken : generate token for user
// return data token (string) or error
func CreateToken(custNumber int) *models.Result {
	accSecret := os.Getenv("ACCESS_SECRET")
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = custNumber
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, _ := at.SignedString([]byte(accSecret))

	return &models.Result{Data: token}
}
