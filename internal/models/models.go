package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	ID          uint   `json:"id" gorm:"primaryyKey"`
	Description string `json:"description" gorm:"text"`
	Note        string `json:"note" gorm:"text"`

	UserID uint `json:"user_id"`
	User   User `json:"user" gorm:"foreignKey:UserID;references:ID"`
}
