package services

import (
	"ProjetFilViolet/backend/api/config"
	"ProjetFilViolet/backend/api/models"
)

func AddFriend(userID uint, friendID uint) error {
	friend := &models.Friend{
		UserID:   userID,
		FriendID: friendID,
	}
	if err := config.DB.Create(friend).Error; err != nil {
		return err
	}
	return nil
}

func RemoveFriend(userID uint, friendID uint) error {
	if err := config.DB.Where("user_id = ? AND friend_id = ?", userID, friendID).Delete(&models.Friend{}).Error; err != nil {
		return err
	}
	return nil
}

func GetFriends(userID uint) ([]models.User, error) {
	var friends []models.User
	if err := config.DB.Table("friends").Select("users.*").Joins("join users on friends.friend_id = users.id").Where("friends.user_id = ?", userID).Scan(&friends).Error; err != nil {
		return nil, err
	}
	return friends, nil
}
