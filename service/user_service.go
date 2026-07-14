package service

import (
	"errors"
	"strings"

	"github.com/EzraArafa/artha-economy-api/model"
	"github.com/EzraArafa/artha-economy-api/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(user *model.User) error {
	// 1. Logika Bisnis: Cegah username kosong
	if strings.TrimSpace(user.Username) == "" {
		return errors.New("username tidak boleh kosong")
	}

	// 2. Logika Bisnis: Validasi struktur hierarki (Role)
	// Ubah huruf menjadi kecil semua agar aman (misal inputnya "President" jadi "president")
	user.Role = strings.ToLower(user.Role)

	validRoles := map[string]bool{
		"president": true,
		"member":    true,
		"prospect":  true,
	}

	if !validRoles[user.Role] {
		return errors.New("role tidak valid. Hanya menerima: president, member, atau prospect")
	}

	// 3. Jika semua logika di atas lolos, perintahkan Repository untuk simpan ke database
	err := s.userRepo.Create(user)
	if err != nil {
		return err
	}

	return nil //mengembalikan nil, berarti proses sukses tanpa error
}

// Fungsi otak dari transaksi
func (s *UserService) TransferBalance(senderID int, receiverID int, amount int) error {
	//Validasi dasar, Uang yang dikirim tidak boleh minus atau 0
	if amount <= 0 {
		return errors.New("jumlah transfer kamu harus lebih dari 0")
	}

	//Mencegah transfer ke diri sendiri
	if senderID == receiverID {
		return errors.New("tidak bisa mentransfer ke diri sendiri")
	}

	//Mencari data pengirim
	sender, err := s.userRepo.FindByID(senderID)
	if err != nil {
		return errors.New("data pengirim tidak ditemukan")
	}

	//Mencari daata penerima
	receiver, err := s.userRepo.FindByID(receiverID)
	if err != nil {
		return errors.New("data penerima tidak ditemukan")
	}

	//Mengecek apakah saldo pengirim cukup
	if sender.Balance < amount {
		return errors.New("saldo tidak cukup untuk melakukan transfer")
	}

	sender.Balance -= amount
	receiver.Balance += amount

	return s.userRepo.ExecuteTransfer(sender, receiver)
}

// Fungsi untuk mengecek syarat sebelum mengeksekusi pembelian
func (s *UserService) PurchaseItem(userID int, item *model.Item, quantity int) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("data user tidak ditemukan")
	}

	totalPrice := item.Price * quantity

	//Validasi saldo
	if user.Balance < totalPrice {
		return errors.New("saldo tidak cukup untuk membeli barang ini")
	}

	return s.userRepo.ExecutePurchase(user, int(item.ID), item.Price, quantity)
}

// Fungsi untuk melihan inventory
func (s *UserService) GetUserInventory(userID int) ([]model.UserInventory, error) {
	return s.userRepo.GetInventoryByUserID(userID)
}

// Fungsi untuk memvalidasi penggunaan barang
func (s *UserService) ConsumeItem(userID int, itemID int, quantity int) error {
	//Mencegah input quantity minus atau nol
	if quantity <= 0 {
		return errors.New("jumlah barang yang digunakan harus lebih dari 0")
	}
	return s.userRepo.ConsumeItem(userID, itemID, quantity)
}

// Fungsi validasi sebelum transfer barang dieksekusi
func (s *UserService) TransferItem(senderID int, receiverID int, itemID int, quantity int) error {
	if senderID == receiverID {
		return errors.New("tidak bisa memberikan barang ke diri sendiri")
	}
	if quantity <= 0 {
		return errors.New("jumlah barang yan diberikan harus lebih dari 0")
	}

	return s.userRepo.TransferItem(senderID, receiverID, itemID, quantity)
}
