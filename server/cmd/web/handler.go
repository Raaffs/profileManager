package main

import (
	"errors"
	"net/http"

	"github.com/Raaffs/profileManager/server/internal/cipher"
	"github.com/Raaffs/profileManager/server/internal/env"
	"github.com/Raaffs/profileManager/server/internal/models"
	"github.com/labstack/echo/v4"
)

func (app *Application) Login(c echo.Context) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&input); err != nil {
		app.logger.Errorf("error binding json to type user \n%w", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	user, err := app.repo.Users.GetUserByEmail(c.Request().Context(), c.FormValue("email"))
	if err != nil {
		if errors.Is(err, models.NotFound) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user not found"})
		}
		app.logger.Errorf("error fetching user by email \n%w", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	token, err := app.GenerateToken(user.ID)
	if err != nil {
		app.logger.Errorf("error generating token \n%w", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}

func (app *Application) Register(c echo.Context) error {
	var u models.User
	if err := c.Bind(&u); err != nil {
		app.logger.Errorf("error binding json to type user \n%w", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	err := app.repo.Users.CreateUser(c.Request().Context(), &u)
	if err != nil {
		app.logger.Errorf("error creating user \n%w", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return nil
}

func (app *Application) CreateProfile(c echo.Context) error {
	var p models.Profile
	if err := c.Bind(&p); err != nil {
		app.logger.Errorf("error binding json to type profile \n%w", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	userID, err := app.GetUserJWT(c)
	if err != nil {
		app.logger.Errorf("error getting user from jwt \n%w", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	p.UserID = userID
	p.AadhaarNumber,err=cipher.Encrypt(p.AadhaarNumber,app.env[env.AES_KEY])
	if err!=nil{
		app.logger.Errorf("error encrypting aadhaar number \n%w", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	if err := app.repo.Profiles.CreateProfile(c.Request().Context(), p); err != nil {
		if errors.Is(err, models.NotFound) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "user not found"})
		}
		app.logger.Errorf("error creating profile \n%w", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	return nil
}

func (app *Application) GetProfile(c echo.Context) error {
	userID, err := app.GetUserJWT(c)
	if err != nil {
		app.logger.Errorf("error getting user from jwt \n%w", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	profile, err := app.repo.Profiles.GetProfileByUserID(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, models.NotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "profile not found"})
		}
		app.logger.Errorf("error fetching profile by user id \n%w", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	profile.AadhaarNumber,err=cipher.Decrypt(profile.AadhaarNumber,app.env[env.AES_KEY])

	return c.JSON(http.StatusOK, profile)
}

func (app *Application) UpdateProfile(c echo.Context) error {
	userID, err := app.GetUserJWT(c)
	if err != nil {
		app.logger.Errorf("error getting user from jwt \n%w", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	var p models.Profile
	if err := c.Bind(&p); err != nil {
		app.logger.Errorf("error binding json to type profile \n%w", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	
	p.UserID = userID
	p.AadhaarNumber,err=cipher.Encrypt(p.AadhaarNumber,app.env[env.AES_KEY])
	if err!=nil{
		app.logger.Errorf("error encrypting aadhaar number \n%w", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	
	if err := app.repo.Profiles.UpdateProfile(c.Request().Context(), p); err != nil {
		if errors.Is(err, models.NotFound) {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "profile not found"})
		}
		app.logger.Errorf("error updating profile \n%w", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}	
	return nil
}
