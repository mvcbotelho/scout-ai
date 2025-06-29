package models

import "gorm.io/gorm"

type Player struct {
	gorm.Model
	Name     string `json:"name" binding:"required" gorm:"not null"`
	Age      int    `json:"age" binding:"required,min=1,max=100" gorm:"not null"`
	Position string `json:"position" binding:"required" gorm:"not null"`
	Team     string `json:"team" binding:"required" gorm:"not null"`
	Goals    int    `json:"goals" binding:"min=0" gorm:"default:0"`
	Tackles  int    `json:"tackles" binding:"min=0" gorm:"default:0"`
	Passes   int    `json:"passes" binding:"min=0" gorm:"default:0"`
}

// TableName especifica o nome da tabela
func (Player) TableName() string {
	return "players"
}
