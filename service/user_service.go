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
