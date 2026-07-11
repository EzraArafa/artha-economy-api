package repository

import (
	"github.com/EzraArafa/artha-economy-api/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	err := r.db.Create(user).Error
	return err
}

// Fungsi untuk mencari data User berdsarkan ID
func (r *UserRepository) FindByID(id int) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	return &user, err
}

// Fungsi untuk mengeksekusi perpindahan uang dengan aman (DB Trannsaction)
func (r *UserRepository) ExecuteTransfer(sender *model.User, receiver *model.User) error {
	tx := r.db.Begin()

	if err := tx.Save(sender).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Save(receiver).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
