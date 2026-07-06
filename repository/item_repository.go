package repository

import (
	"github.com/EzraArafa/artha-economy-api/model"
	"gorm.io/gorm"
)

type ItemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) Create(item *model.Item) error {
	err := r.db.Create(item).Error
	return err
}
