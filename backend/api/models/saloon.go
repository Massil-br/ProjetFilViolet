package models

type Saloon struct {
	Model
	Name        string  `json:"name" gorm:"unique"`
	TalbleCount uint64  `json:"table_count" gorm:"column:tables_count"`
	PlayerCount uint64  `json:"player_count" gorm:"column:players_count"`
	MinBet      uint64  `json:"min_bet"`
	MaxBet      uint64  `json:"max_bet"`
}
