package controller

import (
	"net/http"

	"github.com/EzraArafa/artha-economy-api/model"
	"github.com/EzraArafa/artha-economy-api/service"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService: userService}
}

func (ctrl *UserController) Create(c *gin.Context) {
	var user model.User
	// 1. Tangkap data JSON dari request dan cocokkan dengan cetak biru (struct) User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak sesuai"})
		return
	}

	// 2. Serahkan data ke Service untuk divalidasi logikanya dan disimpan
	err := ctrl.userService.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Jika semua lancar, kembalikan status 201 (Created) dan tampilkan data user yang berhasil dibuat
	c.JSON(http.StatusCreated, gin.H{
		"message": "User berhasil didaftarkan!",
		"data":    user,
	})
}

type TransferInput struct {
	SenderID   int `json:"sender_id"`
	ReceiverID int `json:"receiver_id"`
	Amount     int `json:"amount"`
}

// Fungsi untuk menerima pesanan transfer
func (ctrl *UserController) Transfer(c *gin.Context) {
	var input TransferInput

	//Menangkap format JSON yang dikirim User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak sesuai"})
		return
	}

	//Meneruskan ke Service untuk dieksekusi logikanya
	err := ctrl.userService.TransferBalance(input.SenderID, input.ReceiverID, input.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Jika sukses tanpa error, berikan pesan berhasil
	c.JSON(http.StatusOK, gin.H{
		"message": "Transfer berhasil dilakukan!",
	})

}
