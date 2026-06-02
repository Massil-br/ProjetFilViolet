package services

import (
	"ProjetFilViolet/backend/api/config"
	"ProjetFilViolet/backend/api/models"

	"gorm.io/gorm"
)

func CreateTable(saloonID uint, slotsAvailable uint64) (*models.Table, error) {
	table := &models.Table{
		SaloonID:       saloonID,
		SlotsAvailable: slotsAvailable,
	}
	if err := config.DB.Create(table).Error; err != nil {
		return nil, err
	}
	return table, nil
}

func UpdateTable(table *models.Table) (*models.Table, error) {
	if err := config.DB.Save(table).Error; err != nil {
		return nil, err
	}
	return table, nil
}

func GetTableByID(id uint) (*models.Table, error) {
	var table models.Table
	if err := config.DB.First(&table, id).Error; err != nil {
		return nil, err
	}
	return &table, nil
}

func GetAllTables() ([]models.Table, error) {
	var tables []models.Table
	if err := config.DB.Find(&tables).Error; err != nil {
		return nil, err
	}
	return tables, nil
}

func DeleteTableByID(id uint) error {
	result := config.DB.Delete(&models.Table{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
