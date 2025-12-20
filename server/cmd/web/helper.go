package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Raaffs/profileManager/server/internal/cipher"
	"github.com/Raaffs/profileManager/server/internal/env"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type HttpResponseMsg string

var ErrInternalServer=HttpResponseMsg("internal server error")
var ErrBadRequest=HttpResponseMsg("bad request")
var ErrUnauthorized=HttpResponseMsg("you're not authorized to perform this action")

type JwtCustomClaims struct {
	UserID  int `json:"user_id"`
	jwt.RegisteredClaims
}
func (app *Application) GenerateToken(userID int) (string, error) {
    claims:=&JwtCustomClaims{
		userID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),		},
	}
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(app.env[env.JWT_SECRET]))
}

func (app *Application) GetUserJWT(c echo.Context) (int, error) {
    user, ok := c.Get("user").(*jwt.Token)
    if !ok {
        return 0, errors.New("invalid token type in context")
    }
    claims, ok := user.Claims.(*JwtCustomClaims)
    if !ok {
        return 0, errors.New("invalid token claims")
    }
    userID := claims.UserID
	log.Println("user id: ",userID,claims,user)
    return int(userID), nil
}

func EncryptFields(secretKey string, fields ...*string) error {
    for i, field := range fields {
        // Skip empty optional fields to avoid storing encrypted empty strings
        if field == nil || *field == "" {
            continue
        }
        encryptedValue, err := cipher.Encrypt(secretKey, *field)
        if err != nil {
            return fmt.Errorf("encryption failed for %s: %w", i, err)
        }
        *field = encryptedValue
    }
    return nil
}

func DecryptFields(secretKey string, fields ...*string) error {
    for _, field := range fields {
        if field == nil || *field == "" {
            continue
        }
        decryptedValue, err := cipher.Decrypt(secretKey, *field)
        if err != nil {
            return err
        }
        *field = decryptedValue
    }
    return nil
}