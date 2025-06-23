package models

import "gorm.io/gorm"

type Player struct {
	gorm.Model
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Position string `json:"position"`
	Team     string `json:"team"`
	Goals    int    `json:"goals"`
	Tackles  int    `json:"tackles"`
	Passes   int    `json:"passes"`
}
