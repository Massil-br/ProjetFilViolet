package controllers

import (
	"net/http"
	"strconv"

	"ProjetFilViolet/backend/api/models"
	"ProjetFilViolet/backend/api/services"

	"github.com/labstack/echo/v4"
)

func GetFriends(c echo.Context) error {
	user := c.Get("user").(*models.User)
	friends, err := services.GetFriends(user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get friends")
	}
	return c.JSON(http.StatusOK, friends)
}

func AddFriendById(c echo.Context) error {
	user := c.Get("user").(*models.User)
	friendIDStr , err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid friend ID")
	}
	friendID := uint(friendIDStr)

	err = services.AddFriend(user.ID, friendID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to add friend")
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Friend added successfully",
	})
}

func AddFriendByName(c echo.Context) error {
	user := c.Get("user").(*models.User)
	var FriendName string
	FriendName = c.Param("name")

	friend, err := services.GetUserByName(FriendName)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Friend not found")
	}
	err = services.AddFriend(user.ID, friend.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to add friend")
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Friend added successfully",
	})
}

func RemoveFriend(c echo.Context) error {
	user := c.Get("user").(*models.User)
	friendIDStr , err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid friend ID")
	}
	friendID := uint(friendIDStr)


	err = services.RemoveFriend(user.ID, friendID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to remove friend")
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Friend removed successfully",
	})
}
