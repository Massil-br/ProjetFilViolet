package services

import (
	"gorm.io/gorm"

	"ProjetFilViolet/backend/api/config"
	"ProjetFilViolet/backend/api/models"
	"ProjetFilViolet/backend/api/utils"
)

// CreateUser creates a new user. It handles DB checks, deletes previous soft-deleted
// accounts with the same email, hashes the password and persists the user.
func CreateUser(firstName, lastName, nickName, email, password string) (*models.User, error) {
	// Check existing nickname
	var existingNickname models.User
	if err := config.DB.Where("nick_name = ?", nickName).First(&existingNickname).Error; err == nil {
		return nil, ErrNicknameInUse
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Check existing user (including soft deleted)
	var existing models.User
	if err := config.DB.Unscoped().Where("email = ?", email).First(&existing).Error; err == nil {
		if existing.DeletedAt.Valid {
			if err := config.DB.Unscoped().Delete(&existing).Error; err != nil {
				return nil, ErrFailedRemovePrevious
			}
		} else {
			return nil, ErrEmailInUse
		}
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Hash password
	hashed, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		FirstName:   firstName,
		LastName:    lastName,
		NickName:    nickName,
		Email:       email,
		Password:    hashed,
		Role:        models.UserRole,
		Money:       config.InitialUserMoney,
		IsVerified:  false,
		Status:      models.MainMenu,
		IsConnected: false,
	}

	if err := config.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := config.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func GetUserByID(id string) (*models.User, error) {
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func DeleteUserByID(id string) error {
	result := config.DB.Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func UpdateUser(user *models.User) error {
	if err := config.DB.Save(user).Error; err != nil {
		return err
	}
	return nil
}

// GetUserByEmail returns a user by email.
func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

