package services

import (
	"ProjetFilViolet/backend/api/config"
	"ProjetFilViolet/backend/api/models"
	"errors"
)

func CreateSaloon(name string, playerCount uint64, minBet uint64, maxBet uint64) (*models.Saloon, error) {
	saloon := &models.Saloon{
		Name:        name,
		PlayerCount: playerCount,
		MinBet:      minBet,
		MaxBet:      maxBet,
	}
	if err := config.DB.Create(saloon).Error; err != nil {
		return nil, err
	}
	return saloon, nil
}


func UpdateSaloon(Id uint, Name string, PlayerCount uint64, MinBet uint64, MaxBet uint64) (*models.Saloon, error) {
	saloon,err  := GetSaloonByID(Id)
	if err != nil {
		return nil, err
	}
	
	if MaxBet < MinBet {
		return nil, errors.New("Max bet must be higher than Min bet")
	}

	saloon.Name = Name
	saloon.PlayerCount = PlayerCount
	saloon.MinBet = MinBet
	saloon.MaxBet = MaxBet
	if err := config.DB.Save(saloon).Error; err != nil {
		return nil, err
	}
	return saloon, nil
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