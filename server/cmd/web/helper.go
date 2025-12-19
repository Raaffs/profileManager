package main

import (
	"errors"
	"time"

	"github.com/Raaffs/profileManager/server/internal/env"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)
func (app *Application) GenerateToken(userID int) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(time.Hour * 72).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(app.env[env.JWT_SECRET]))
}

func (app *Application) GetUserJWT(c echo.Context) (int, error) {
	user:=c.Get("user").(*jwt.Token)
	claims:=user.Claims.(jwt.MapClaims)
	userID,ok:=claims["user_id"].(float64)
	if !ok{
		return 0,errors.New("invalid token")
	}
	return int(userID),nil
}