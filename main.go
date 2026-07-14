package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	// Pastikan mengubah <username-github-mu> sesuai milikmu
	"github.com/EzraArafa/artha-economy-api/config"
	"github.com/EzraArafa/artha-economy-api/controller"
	"github.com/EzraArafa/artha-economy-api/model"
	"github.com/EzraArafa/artha-economy-api/repository"
	"github.com/EzraArafa/artha-economy-api/service"
)

func main() {
	// 1. Koneksi Database
	config.ConnectDatabase()

	// 2. AutoMigrate Tabel
	err := config.DB.AutoMigrate(&model.User{}, &model.Item{}, &model.UserInventory{})
	if err != nil {
		fmt.Println("Gagal migrasi:", err)
	}

	//===KOMPONEN ITEM
	itemRepo := repository.NewItemRepository(config.DB)
	itemService := service.NewItemService(itemRepo)
	itemController := controller.NewItemController(itemService)

	//===KOMPONEN USER
	userRepo := repository.NewUserRepository(config.DB)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService, itemService)

	// 4. Inisialisasi Router
	r := gin.Default()

	// 5. Mendaftarkan Jalur (Routes)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Artha Economy API berjalan!"})
	})

	// Ini adalah rute baru kita. Jika ada request POST ke /users, serahkan ke Pelayan (Controller)
	r.POST("/users", userController.Create)
	r.POST("/item", itemController.Create)
	r.POST("/buy", userController.BuyItem)
	r.POST("/transfer", userController.Transfer)

	r.GET("/users/:user_id/inventory", userController.GetInventory)

	r.POST("/use-item", userController.ConsumeItem)
	r.POST("/give-item", userController.TransferItem)
	// 6. Jalankan Server
	r.Run(":8080")
}
