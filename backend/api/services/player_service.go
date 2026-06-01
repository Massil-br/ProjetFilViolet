package services

import (
	"ProjetFilViolet/backend/api/config"
	"ProjetFilViolet/backend/api/models"

	"gorm.io/gorm"
)

func CreatePlayer(player *models.Player) error {
	if err := config.DB.Create(player).Error; err != nil {
		return err
	}
	return nil
}

func GetPlayerByID(id uint) (*models.Player, error) {
	var player models.Player	
	if err := config.DB.First(&player, id).Error; err != nil {
		return nil, err
	}
	return &player, nil
}

func GetPlayerByName(name string) (*models.Player, error) {
	var player models.Player
	if err := config.DB.Joins("JOIN \"user\" ON \"user\".id = player.user_id").Where("\"user\".nick_name = ?", name).First(&player).Error; err != nil {
		return nil, err
	}
	return &player, nil
}

func GetAllPlayersByTableId(tableID uint) ([]models.Player, error) {
	var players []models.Player
	if err := config.DB.Table("players").Select("players.*").Joins("join player_tables on player_tables.player_id = players.id").Where("player_tables.table_id = ?", tableID).Scan(&players).Error; err != nil {
		return nil, err
	}
	return players, nil
}

func DeletePlayerByID(id uint) error {
	result := config.DB.Delete(&models.Player{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func UpdatePlayer(player *models.Player) error {
	if err := config.DB.Save(player).Error; err != nil {
		return err
	}	
	return nil
}


func SetPlayerInitialBet(playerID uint, bet uint64) error {
	var player models.Player
	if err := config.DB.First(&player, playerID).Error; err != nil {
		return err
	}	
	player.InitialBet = bet
	if err := config.DB.Save(&player).Error; err != nil {
		return err
	}
	return nil
}

func AddMoneyToPlayer(playerID uint, amount uint64) error {
	var player models.Player	
	if err := config.DB.First(&player, playerID).Error; err != nil {
		return err
	}	
	player.Money += amount
	if err := config.DB.Save(&player).Error; err != nil {
		return err
	}
	return nil
}

func RemoveMoneyFromPlayer(playerID uint, amount uint64) error{
	var player models.Player
	if err := config.DB.First(&player, playerID).Error; err != nil {
		return err	
	}
	if player.Money < amount {
		return gorm.ErrInvalidData
	}
	player.Money -= amount
	if err := config.DB.Save(&player).Error; err != nil {
		return err
	}
	return nil
}
