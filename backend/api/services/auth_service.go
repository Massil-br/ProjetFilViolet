package services

import (
	"errors"
	"os"
	"time"

	"ProjetFilViolet/backend/api/config"
	"ProjetFilViolet/backend/api/models"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrEmailInUse           = errors.New("email already in use")
	ErrNicknameInUse        = errors.New("nickname already in use")
	ErrFailedRemovePrevious = errors.New("failed to remove previous deleted user")
)



func ConnectUser(userID uint) error {
	if err := config.DB.Model(&models.User{}).Where("id = ?", userID).Update("is_connected", true).Error; err != nil {
		return err
	}
	return nil
}

func DisconnectUser(userID uint) error {
	if err := config.DB.Model(&models.User{}).Where("id = ?", userID).Update("is_connected", false).Error; err != nil {
		return err
	}
	return nil
}

// GenerateTokenForUserID returns a signed JWT containing only the user id.
func GenerateTokenForUserID(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	secret := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(secret))
}

