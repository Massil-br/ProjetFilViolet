package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func MainPage(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Welcome to the API main page.",
	})
}
