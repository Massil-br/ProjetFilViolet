package services

import (
	"time"

	"ProjetFilViolet/backend/api/config"
	"ProjetFilViolet/backend/api/models"

	"gorm.io/gorm"
)

func CreateTchatMessage(UserID uint, TableID uint, Content string) error {
	tchatMessage := &models.TchatMessage{
		UserID:  UserID,
		TableID: TableID,
		Message: Content,
	}
	if err := config.DB.Create(tchatMessage).Error; err != nil {
		return err
	}
	return nil
}

func DeleteTchatMessageByID(id uint) error {
	result := config.DB.Delete(&models.TchatMessage{}, id)
	if result.Error != nil {
		return result.Error	
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil	
}

func GetTchatMessagesByTableID(tableID uint) ([]models.TchatMessage, error){
	var tchatMessages []models.TchatMessage
	tenMinutesAgo := time.Now().Add(-10 * time.Minute)
	if err := config.DB.Where("table_id = ? AND created_at >= ?", tableID, tenMinutesAgo).Find(&tchatMessages).Error; err != nil {
		return nil, err
	}
	return tchatMessages, nil
}