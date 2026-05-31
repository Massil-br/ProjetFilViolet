package models

import (
	"time"

	"gorm.io/gorm"
)

type PlayerStatus int

const (
	MainMenu PlayerStatus = iota
	InQueue
	Ingame
)

type Model struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Role int

const (
	UserRole Role = iota
	Moderator
	Admin
)

type User struct {
	Model
	FirstName   string       `json:"first_name"`
	LastName    string       `json:"last_name"`
	NickName    string       `json:"nick_name" gorm:"unique"`
	Email       string       `json:"email" gorm:"unique"`
	Password    string       `json:"password"`
	IsVerified  bool         `json:"is_verified"`
	Status      PlayerStatus `json:"status" gorm:"type:VARCHAR(255);default:MainMenu"`
	IsConnected bool         `json:"is_connected"`
	Money       uint64       `json:"money"`
	Role        Role         `json:"role" gorm:"type:VARCHAR(255);default:user"`
}
