package models

type Player struct {
	UserID     uint    `json:"user_id" gorm:"primaryKey;foreignKey:UserID;references:ID;unique"`
	TableID    uint    `json:"table_id" gorm:"primaryKey;foreignKey:TableID;references:ID"`
	InitialBet float64 `json:"initial_bet"`
	Money      float64 `json:"money"`
}
