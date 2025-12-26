package main

import (
	"github.com/labstack/echo/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
)
func (app *Application) Authenticate() echo.MiddlewareFunc {
    return echojwt.WithConfig(echojwt.Config{
        SigningKey: []byte(app.env["JWT_SECRET"]),
        TokenLookup: "header:Authorization:Bearer ",
    })
}