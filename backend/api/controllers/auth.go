package controllers

import (
	"errors"
	"net/http"
	"regexp"

	"ProjetFilViolet/backend/api/models"
	"ProjetFilViolet/backend/api/services"
	"ProjetFilViolet/backend/api/utils"

	"github.com/labstack/echo/v4"
)

type createUserReq struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	NickName        string `json:"nick_name"`
	Email           string `json:"email"`
	ConfirmEmail	 string `json:"confirmEmail"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUser(c echo.Context) error {
	var req createUserReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	if req.ConfirmEmail != req.Email{
		return c.JSON(http.StatusBadRequest, echo.Map{"error":"email and confirmEmail do not match"})
	}

	// Validation in controller
	if req.ConfirmPassword != req.Password {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "passwords do not match"})
	}
	if len(req.Password) < 8 || !regexp.MustCompile("[0-9]").MatchString(req.Password) {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "password must be at least 8 chars and contain a digit"})
	}

	// Delegate creation to service (service handles hashing and DB checks)
	user, err := services.CreateUser(req.FirstName, req.LastName, req.NickName, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrEmailInUse) {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "email already in use"})
		}
		if errors.Is(err, services.ErrNicknameInUse) {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "nickname already in use"})
		}
		if errors.Is(err, services.ErrFailedRemovePrevious) {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to remove previous deleted user"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not create user {" + err.Error() + "}"})
	}

	// mask password before returning
	user.Password = ""

	return c.JSON(http.StatusCreated, echo.Map{"message": "user created", "user": user})
}

func Login(c echo.Context) error {
	var req loginReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid input"})
	}

	user, err := services.GetUserByEmail(req.Email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid credentials"})
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid credentials"})
	}

	token, err := services.GenerateTokenForUserID(user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not generate token"})
	}

	user.Password = ""
	if err := services.ConnectUser(user.ID); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not update user connection status"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "login successful", "token": token, "user": user})
}

func Logout(c echo.Context) error {
	// Extract user from context (set by AuthMiddleware)
	user, ok := c.Get("user").(*models.User)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid or missing token"})
	}

	// Disconnect user
	if err := services.DisconnectUser(user.ID); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "could not logout"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "logout successful"})
}
