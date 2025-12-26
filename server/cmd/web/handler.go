package main

import (
	"errors"
	"net/http"

	"github.com/Raaffs/profileManager/server/internal/cipher"
	"github.com/Raaffs/profileManager/server/internal/env"
	"github.com/Raaffs/profileManager/server/internal/models"
	"github.com/Raaffs/profileManager/server/internal/utils"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (app *Application) Login(c echo.Context) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&input); err != nil {
		app.logger.Errorf("error binding json to type user \n%w", err)
		return c.JSON(http.StatusBadRequest, map[string]HttpResponseMsg{"error": ErrBadRequest})
	}

	user, err := app.repo.Users.GetByEmail(c.Request().Context(), input.Email)
	if err != nil {
		if errors.Is(err, models.NotFound) {
			return c.JSON(http.StatusNotFound, map[string]HttpResponseMsg{"error": ErrNotFound})
		}
		app.logger.Errorf("error fetching user by email \n%w", err)
		app.health.SetStatus(StatusDegraded)
		return c.JSON(http.StatusInternalServerError, map[string]HttpResponseMsg{"error": ErrInternalServer})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return c.JSON(http.StatusUnauthorized, map[string]HttpResponseMsg{"error": "invalid username or password"})
		}
		app.health.SetStatus(StatusDegraded)
		app.logger.Error("error comparing password hash \n%w", err)
		return c.JSON(http.StatusInternalServerError, map[string]HttpResponseMsg{"error": ErrInternalServer})
	}

	token, err := app.GenerateToken(user.ID)
	if err != nil {
		app.health.SetStatus(StatusCritical)
		app.logger.Errorf("error generating token \n%w", err)
		return c.JSON(http.StatusInternalServerError, map[string]HttpResponseMsg{"error": ErrInternalServer})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}

func (app *Application) Register(c echo.Context) error {
	var u struct{
		Email string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
	}
	if err := c.Bind(&u); err != nil {
		app.logger.Errorf("error binding json to type user \n%w", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	
	validate := utils.NewValidator()
	validate.NameLength(u.Username, 3, 20)
	validate.Mail(u.Email)

	if !validate.Valid() {
		return c.JSON(http.StatusBadRequest, validate.Errors)
	}

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		app.health.SetStatus(StatusDegraded)
		app.logger.Errorf("error hashing password \n%w", err)
		return c.JSON(http.StatusInternalServerError, map[string]HttpResponseMsg{"error": ErrInternalServer})
	}

	var user models.User

	user.Email=u.Email
	user.Username=u.Username
	user.PasswordHash=hashedPassword

	if err := app.repo.Users.Create(c.Request().Context(), &user); err != nil {
		if errors.Is(err, models.AlreadyExists) {
			return c.JSON(http.StatusConflict, map[string]string{"error": "email or username already exists"})
		}
		app.health.SetStatus(StatusDegraded)
		app.logger.Errorf("error creating user \n%w", err)
		return c.JSON(http.StatusInternalServerError, map[string]HttpResponseMsg{"error": ErrInternalServer})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "account created successfully"})
}

func (app *Application) CreateProfile(c echo.Context) error {
	var p models.Profile
	if err := c.Bind(&p); err != nil {
		app.logger.Errorf("error binding json to type profile \n%w", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	userID, err := app.GetUserJWT(c)
	if err != nil {
		if errors.Is(err, ErrInvalidToken){
			return c.JSON(http.StatusUnauthorized, map[string]HttpResponseMsg{"error": ErrUnauthorized})
		}
		app.health.SetStatus(StatusCritical)
		app.logger.Errorf("error getting user from jwt \n%w", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}
	if validate := ValidateProfile(p); !validate.Valid(){
		return c.JSON(http.StatusBadRequest, validate.Errors)
	}

	p.UserID = userID

	if err := EncryptFields(app.env[env.AES_KEY],&p.AadhaarNumber); err!=nil{
		app.health.SetStatus(StatusCritical)
		app.logger.Errorf("CRITICAL ERROR : cipher failure \n%w", err)	
		return c.JSON(http.StatusInternalServerError, map[string]HttpResponseMsg{"error": ErrInternalServer})
	}

	if err := app.repo.Profiles.Create(c.Request().Context(), p); err != nil {
		if errors.Is(err, models.NotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
		}

		if errors.Is(err, models.AlreadyExists) {
			return c.JSON(http.StatusConflict, map[string]string{"error": "phone no. already exists"})
		}
		
		app.health.SetStatus(StatusDegraded)
		app.logger.Errorf("error creating profile \n%w", err)
		return c.JSON(http.StatusInternalServerError, map[string]HttpResponseMsg{"error": ErrInternalServer})
	}

	return c.JSON(http.StatusOK,map[string]string{"message":"profile created successfully"})
}

func (app *Application) GetProfile(c echo.Context) error {
	userID, err := app.GetUserJWT(c)
	if err != nil {
		app.logger.Errorf("error getting user from jwt \n%w", err)
		return c.JSON(http.StatusUnauthorized, map[string]HttpResponseMsg{"error": ErrUnauthorized})
	}

	profile, err := app.repo.Profiles.GetByUserID(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, models.NotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "profile not found"})
		}
		app.logger.Errorf("error fetching profile by user id \n%w", err)
		return c.JSON(http.StatusInternalServerError, map[string]HttpResponseMsg{"error": ErrInternalServer})
	}

	if err := DecryptFields(app.env[env.AES_KEY], &profile.AadhaarNumber);err != nil {
		app.logger.Errorf("CRITICAL ERROR: cipher failure: \n%w", err)
		app.health.SetStatus(StatusCritical)
		return c.JSON(http.StatusInternalServerError, map[string]HttpResponseMsg{"error": ErrInternalServer})
	}
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
		return c.JSON(http.StatusBadRequest, map[string]HttpResponseMsg{"error": ErrBadRequest})
	}

	if validate := ValidateProfile(p); !validate.Valid(){
		return c.JSON(http.StatusBadRequest, validate.Errors)
	}

	p.UserID = userID
	p.AadhaarNumber, err = cipher.Encrypt(app.env[env.AES_KEY],p.AadhaarNumber)
	if err != nil {
		app.health.SetStatus(StatusCritical)
		app.logger.Errorf("CRITICAL ERROR: cipher failure \n%w", err)
		return c.JSON(http.StatusInternalServerError, map[string]HttpResponseMsg{"error": ErrInternalServer})
	}

	if err := app.repo.Profiles.Update(c.Request().Context(), p); err != nil {
		if errors.Is(err, models.NotFound) {
			return c.JSON(http.StatusBadRequest, map[string]HttpResponseMsg{"error": ErrNotFound})
		}
		if errors.Is(err, models.AlreadyExists) {
			//the phone number is the only unique field that can cause conflict here
			//that's why we return this specific message
			return c.JSON(http.StatusConflict, map[string]string{"error": "phone no. already exists"})
		}
		app.health.SetStatus(StatusDegraded)
		app.logger.Errorf("error updating profile \n%w", err)
		return c.JSON(http.StatusInternalServerError, map[string]HttpResponseMsg{"error": ErrInternalServer})
	}
	return c.JSON(http.StatusOK,map[string]string{
		"message":"profile updated successfully",
	})
}

