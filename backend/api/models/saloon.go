package models

type Saloon struct {
	Model
	TalbleCount uint64  `json:"table_count" gorm:"column:tables_count"`
	PlayerCount uint64  `json:"player_count" gorm:"column:players_count"`
	MinBet      float64 `json:"min_bet"`
	MaxBet      float64 `json:"max_bet"`
}
