package repository

import (
	"errors"

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

// Fungsi untuk memakai barang
func (r *UserRepository) ConsumeItem(userID int, itemID int, quantity int) error {
	var inventory model.UserInventory

	//Cari barang didalam tas user
	err := r.db.Where("user_id = ? AND item_id = ?", userID, itemID).First(&inventory).Error
	if err != nil {
		return errors.New("barang tidak ditemukan didalam inventory")
	}

	//Cek apakah jumlah barang yang mau dipakai tidak melebihi yang dimiliki
	if inventory.Quantity < quantity {
		return errors.New("jumlah barang tidak mencukupi untuk digunakan")
	}

	//Kurangi jumlah
	inventory.Quantity -= quantity

	//Jika habis (0), hapus dari tas. Jika masih ada, simpan pembaruannya
	if inventory.Quantity == 0 {
		return r.db.Delete(&inventory).Error
	}

	return r.db.Save(&inventory).Error
}

// Fungsi untuk memindahkan barang antara user dengan aman
func (r *UserRepository) TransferItem(senderID int, receiverID int, itemID int, quantity int) error {
	tx := r.db.Begin()

	//Mengurangi barang dari tas pengirim
	var senderInv model.UserInventory
	if err := tx.Where("user_id = ? AND item_id = ?", senderID, itemID).First(&senderInv).Error; err != nil {
		tx.Rollback()
		return errors.New("pengirim tidak memiliki barang tersebut ditasnya")
	}

	if senderInv.Quantity < quantity {
		tx.Rollback()
		return errors.New("jumlah barang pengirim tidak mencukupi untuk diberikan")
	}

	senderInv.Quantity -= quantity
	if senderInv.Quantity == 0 {
		if err := tx.Delete(&senderInv).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		if err := tx.Save(&senderInv).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	//Menambah barang ke tas penerima
	var receiverInv model.UserInventory
	err := tx.Where("user_id = ? AND item_id = ?", receiverID, itemID).First(&receiverInv).Error

	if err != nil {
		//Jika belum punya barang, maka membuat tumpukan baru
		receiverInv = model.UserInventory{
			UserID:   receiverID,
			ItemID:   itemID,
			Quantity: quantity,
		}
		if err := tx.Create(&receiverInv).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		//Jika sudah punya, tambahkan ke tumpukan yang ada
		receiverInv.Quantity += quantity
		if err = tx.Save(&receiverInv).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
