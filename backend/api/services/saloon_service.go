package services

import (
	"ProjetFilViolet/backend/api/config"
	"ProjetFilViolet/backend/api/models"
)

func CreateSaloon(saloon *models.Saloon) error {
	if err := config.DB.Create(saloon).Error; err != nil {
		return err
	}
	return nil
}

func UpdateSaloon(saloon *models.Saloon) error {
	if err := config.DB.Save(saloon).Error; err != nil {
		return err
	}
	return nil
}



func GetSaloonByID(id uint) (*models.Saloon, error) {
	var saloon models.Saloon
	if err := config.DB.First(&saloon, id).Error; err != nil {
		return nil, err
	}
	return &saloon, nil
}

func GetSaloonByName(name string) (*models.Saloon, error) {
	var saloon models.Saloon
	if err := config.DB.Where("name = ?", name).First(&saloon).Error; err != nil {
		return nil, err
	}
	return &saloon, nil
}

func ResearchSaloonByName(name string) ([]models.Saloon, error) {
	var saloons []models.Saloon
	if err := config.DB.Where("name ILIKE ?", "%"+name+"%").Find(&saloons).Error; err != nil {
		return nil, err
	}	
	return saloons, nil
}

func GetAllSaloons() ([]models.Saloon, error) {
	var saloons []models.Saloon
	if err := config.DB.Find(&saloons).Error; err != nil {
		return nil, err
	}
	return saloons, nil
}

func DeleteSaloonByID(id uint) error {
	result := config.DB.Delete(&models.Saloon{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}


func DeleteSaloonByName(name string) error {
	result := config.DB.Where("name = ?", name).Delete(&models.Saloon{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}