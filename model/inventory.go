package model

import "time"

type UserInventory struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	UserID    int       `json:"user_id"`
	ItemID    int       `json:"item_id"`
	Item      Item      `json:"item" gorm:"foreignKey:ItemID"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
