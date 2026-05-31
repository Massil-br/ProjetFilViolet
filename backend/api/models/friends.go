package models

type Friend struct {
	UserID   uint `json:"user_id" gorm:"primaryKey;foreignKey:UserID;references:ID"`
	FriendID uint `json:"friend_id" gorm:"primaryKey;foreignKey:FriendID;references:ID"`
}
