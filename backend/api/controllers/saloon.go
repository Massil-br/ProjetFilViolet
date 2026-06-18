package controllers

import (
	"ProjetFilViolet/backend/api/services"
	"strconv"

	"github.com/labstack/echo/v4"
)

type SaloonRequest struct {
	Name        string `json:"name"`
	PlayerCount uint64 `json:"player_count"`
	MinBet      uint64 `json:"min_bet"`
	MaxBet      uint64 `json:"max_bet"`
}

type UpdateSaloonRequest struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	PlayerCount uint64 `json:"player_count"`
	MinBet      uint64 `json:"min_bet"`
	MaxBet      uint64 `json:"max_bet"`
}

func CreateSaloon(c echo.Context) error {
	var req SaloonRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(400, "Invalid request body")
	}
	saloon, err := services.CreateSaloon(req.Name, req.PlayerCount, req.MinBet, req.MaxBet)
	if err != nil {
		return echo.NewHTTPError(500, "Failed to create saloon")
	}
	return c.JSON(201, echo.Map{"message": "Saloon created successfully", "saloon": saloon})
}

func UpdateSaloon(c echo.Context) error {
	var req UpdateSaloonRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(400, "Invalid request body")
	}
	saloon, err := services.UpdateSaloon(req.Id, req.Name, req.PlayerCount, req.MinBet, req.MaxBet)
	if err != nil {
		return echo.NewHTTPError(500, "Failed to update saloon")
	}
	return c.JSON(200, echo.Map{"message": "Saloon updated successfully", "saloon": saloon})
}

func GetSaloonByID(c echo.Context) error {
	paramId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(400, "Invalid saloon ID")
	}
	id := uint(paramId)

	saloon, err := services.GetSaloonByID(id)
	if err != nil {
		return echo.NewHTTPError(404, "Saloon not found")
	}
	return c.JSON(200, saloon)
}

func GetSaloonByName(c echo.Context) error {
	name := c.Param("name")
	saloon, err := services.GetSaloonByName(name)
	if err != nil {
		return echo.NewHTTPError(404, "Saloon not found")
	}
	return c.JSON(200, saloon)
}

func GetAllSaloons(c echo.Context) error {
	saloons, err := services.GetAllSaloons()
	if err != nil {
		return echo.NewHTTPError(500, "Failed to fetch saloons")
	}
	return c.JSON(200, saloons)
}


func DeleteSaloonByID(c echo.Context) error {
	paramId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.NewHTTPError(400, "Invalid saloon ID")
	}
	id := uint(paramId)

	err = services.DeleteSaloonByID(id)
	if err != nil {
		return echo.NewHTTPError(500, "Failed to delete saloon")
	}
	return c.JSON(200, echo.Map{"message": "Saloon deleted successfully"})
}

func DeleteSaloonByName(c echo.Context) error {
	name := c.Param("name")
	err := services.DeleteSaloonByName(name)
	if err != nil {
		return echo.NewHTTPError(500, "Failed to delete saloon")
	}
	return c.JSON(200, echo.Map{"message": "Saloon deleted successfully"})
}
