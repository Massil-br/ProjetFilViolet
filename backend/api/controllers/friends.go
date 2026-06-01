package controllers

import (
	"net/http"

	"ProjetFilViolet/backend/api/models"
	"ProjetFilViolet/backend/api/services"

	"github.com/labstack/echo/v4"
)

type FriendByIdRequest struct {
	FriendID uint `json:"friend_id"`
}

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
	var req FriendByIdRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	err := services.AddFriend(user.ID, req.FriendID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to add friend")
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Friend added successfully",
	})
}


type AddFriendByNameRequest struct {
	FriendName string `json:"friend_name"`
}

func AddFriendByName(c echo.Context) error{
	user := c.Get("user").(*models.User)
	var req AddFriendByNameRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	friend, err := services.GetUserByName(req.FriendName)
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
	var req FriendByIdRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}	
	err := services.RemoveFriend(user.ID, req.FriendID)	
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to remove friend")
	}	
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Friend removed successfully",
	})
}
	
