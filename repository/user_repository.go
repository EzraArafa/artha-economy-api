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

// Fungsi untuk menyimpan pembaruan data user tunggal
func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Fungsi untuk memotong saldo dan menambah barang ke inventory
func (r *UserRepository) ExecutePurchase(user *model.User, itemID int, price int, quantity int) error {
	tx := r.db.Begin()

	//Memotong saldo user dan simpan
	user.Balance -= (price * quantity)
	if err := tx.Save(user).Error; err != nil {
		tx.Rollback() //batal jika gagal
		return err
	}

	//Cek apakah barang sudah masuk ke inventory user
	var inventory model.UserInventory
	err := tx.Where("user_id = ? AND item_id = ?", user.ID, itemID).First(&inventory).Error

	if err != nil {
		//Jika barang belum ada, maka membuat tumpukan baru
		inventory = model.UserInventory{
			UserID:   int(user.ID),
			ItemID:   itemID,
			Quantity: quantity,
		}
		if err := tx.Create(&inventory).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		//Jika barang sudah ada
		inventory.Quantity += quantity
		if err := tx.Save(&inventory).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	//Jika semua berhasil
	return tx.Commit().Error
}

// Fungsi untuk mengambil seluruh inventory milik user
func (r *UserRepository) GetInventoryByUserID(userID int) ([]model.UserInventory, error) {
	var inventories []model.UserInventory
	//Preload ("item") akan otomatis memuat semua detail barang
	err := r.db.Preload("Item").Where("user_id = ?", userID).Find(&inventories).Error
	return inventories, err
}
