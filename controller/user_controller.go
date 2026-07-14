package controller

import (
	"net/http"
	"strconv"

	"github.com/EzraArafa/artha-economy-api/model"
	"github.com/EzraArafa/artha-economy-api/service"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *service.UserService
	itemService *service.ItemService
}

func NewUserController(userService *service.UserService, itemService *service.ItemService) *UserController {
	return &UserController{userService: userService, itemService: itemService}
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

// Struktur untuk menangkap data pembelian
type BuyInput struct {
	UserID   int `json:"user_id"`
	ItemID   int `json:"item_id"`
	Quantity int `json:"quantity"`
}

func (ctrl *UserController) BuyItem(c *gin.Context) {
	var input BuyInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak sesuai"})
		return
	}

	//Mencegah pembelian 0 atau minus
	if input.Quantity <= 0 {
		input.Quantity = 1
	}

	//Cek harga barang ke ItemService
	item, err := ctrl.itemService.GetItemByID(input.ItemID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Barang tidak ditemukan"})
		return
	}

	//Menyuruh Service mengeksekusi pembelian (potong saldo & masuk ke inventory)
	err = ctrl.userService.PurchaseItem(input.UserID, item, input.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//memberi respon sukses
	c.JSON(http.StatusOK, gin.H{
		"message":  "Pembelian sukses!",
		"item":     item.Name,
		"price":    item.Price,
		"quantity": input.Quantity,
		"total":    item.Price * input.Quantity,
	})
}

// Fungsi untuk menampilkan inventory user
// Fungsi untuk menampilkan isi tas user
func (ctrl *UserController) GetInventory(c *gin.Context) {
	// 1. Ambil ID dari URL (misal: /users/2/inventory)
	userIDStr := c.Param("user_id")

	// 2. Ubah format ID dari teks ke angka
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format ID user tidak valid"})
		return
	}

	// 3. Minta data ke Service
	inventories, err := ctrl.userService.GetUserInventory(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. Tampilkan data ke layar!
	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil memuat isi tas",
		"data":    inventories,
	})
}

// Struct penampung data request pemakaian barang
type ConsumeInput struct {
	UserID   int `json:"user_id"`
	ItemID   int `json:"item_id"`
	Quantity int `json:"quantity"`
}

// Fungsi memakai barang
func (ctrl *UserController) ConsumeItem(c *gin.Context) {
	var input ConsumeInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak sesuai"})
		return
	}

	// Suruh service mengeksekusi penggunaan barang
	err := ctrl.userService.ConsumeItem(input.UserID, input.ItemID, input.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Barang berhasil digunakan",
		"item_id":  input.ItemID,
		"quantity": input.Quantity,
	})
}

// Kirim barang
type TrasnferItemInput struct {
	SenderID   int `json:"sender_id"`
	ReceiverID int `json:"receiver_id"`
	ItemID     int `json:"item_id"`
	Quantity   int `json:"quantity"`
}

// Fungsi untuk trasnfer barang antar user
func (ctrl *UserController) TransferItem(c *gin.Context) {
	var input TrasnferItemInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak sesuai"})
		return
	}

	// Eksekusi trasnfer melalui Service
	err := ctrl.userService.TransferItem(input.SenderID, input.ReceiverID, input.ItemID, input.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Barang berhasil diserahkan!",
		"item_id":  input.ItemID,
		"quantity": input.Quantity,
	})
}
