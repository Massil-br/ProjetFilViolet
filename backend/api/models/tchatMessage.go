package models

type TchatMessage struct {
	Model
	UserID  uint   `json:"user_id" gorm:"foreignKey:UserID;references:ID"`
	Message string `json:"message"`
	TableID uint   `json:"table_id" gorm:"foreignKey:TableID;references:ID"`
}
