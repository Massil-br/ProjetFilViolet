package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"ProjetFilViolet/backend/api/config"
	"ProjetFilViolet/backend/api/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(minRole models.Role) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Authentification (JWT)
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing or invalid Authorization header")
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			secret := os.Getenv("JWT_SECRET")
			token, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims")
			}

			userIDFloat, ok := claims["user_id"].(float64)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user_id in token")
			}
			userID := uint(userIDFloat)

			user := &models.User{}
			if err := config.DB.First(user, userID).Error; err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "User not found")
			}

			// Vérifier que l'utilisateur est connecté
			if !user.IsConnected {
				return echo.NewHTTPError(http.StatusUnauthorized, "User is logged out")
			}

			// Autorisation (rôle) — comparaison numérique basée sur models.Role
			if user.Role < minRole {
				return echo.ErrForbidden
			}

			c.Set("user", user)

			return next(c)
		}
	}
}
