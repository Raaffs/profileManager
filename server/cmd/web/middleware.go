package main

import (
	"github.com/labstack/echo/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
)
func (app *Application) Authenticate() echo.MiddlewareFunc {
    return echojwt.WithConfig(echojwt.Config{
        // Tell Echo which secret to use to verify tokens
        SigningKey: []byte(app.env["JWT_SECRET"]),
        // This ensures the token is extracted from "Authorization: Bearer <token>"
        TokenLookup: "header:Authorization:Bearer ",
    })
}