package services

import (
	"ProjetFilViolet/backend/api/config"
	"ProjetFilViolet/backend/api/models"

	"gorm.io/gorm"
)



func CreatePlayer(userID uint, initialBet uint64, TableID uint) (*models.Player, error) {
	player := &models.Player{
		UserID: userID,
		InitialBet: initialBet,
		TableID: TableID,
		Money: initialBet,
	}
	if err := config.DB.Create(player).Error; err != nil {
		return nil, err
	}
	return player, nil
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

func UpdatePlayer(player *models.Player) (*models.Player, error) {
	if err := config.DB.Save(player).Error; err != nil {
		return nil, err
	}
	return player, nil
}

func SetPlayerMoney(playerID uint, money uint64) (*models.Player, error) {
	var player models.Player
	if err := config.DB.First(&player, playerID).Error; err != nil {
		return nil, err
	}
	player.Money = money
	if err := config.DB.Save(&player).Error; err != nil {
		return nil, err
	}
	return &player, nil
}



func SetPlayerInitialBet(playerID uint, bet uint64) (*models.Player, error) {
	var player models.Player
	if err := config.DB.First(&player, playerID).Error; err != nil {
		return nil, err
	}	
	player.InitialBet = bet
	if err := config.DB.Save(&player).Error; err != nil {
		return nil, err
	}
	return &player, nil
}

func AddMoneyToPlayer(playerID uint, amount uint64) (*models.Player, error) {
	var player models.Player	
	if err := config.DB.First(&player, playerID).Error; err != nil {
		return nil, err
	}	
	player.Money += amount
	if err := config.DB.Save(&player).Error; err != nil {
		return nil,	 err
	}
	return &player, nil
}

func RemoveMoneyFromPlayer(playerID uint, amount uint64) (*models.Player, error) {
	var player models.Player
	if err := config.DB.First(&player, playerID).Error; err != nil {
		return nil, err
	}
	if player.Money < amount {
		return nil, gorm.ErrInvalidData
	}
	player.Money -= amount
	if err := config.DB.Save(&player).Error; err != nil {
		return nil, err
	}
	return &player, nil
}


