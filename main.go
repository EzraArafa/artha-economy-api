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
	err := config.DB.AutoMigrate(&model.User{}, &model.Item{})
	if err != nil {
		fmt.Println("Gagal migrasi:", err)
	}

	// 3. Merangkai Komponen (Dependency Injection)
	// a. Buat Petugas Gudang dengan memberikan alat kerjanya (Database)
	userRepo := repository.NewUserRepository(config.DB)
	// b. Buat Koki Dapur dengan memberikan asistennya (Petugas Gudang)
	userService := service.NewUserService(userRepo)
	// c. Buat Pelayan dengan memberikan koki andalannya (Service)
	userController := controller.NewUserController(userService)

	// 4. Inisialisasi Router
	r := gin.Default()

	// 5. Mendaftarkan Jalur (Routes)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Artha Economy API berjalan!"})
	})

	// Ini adalah rute baru kita. Jika ada request POST ke /users, serahkan ke Pelayan (Controller)
	r.POST("/users", userController.Create)

	// 6. Jalankan Server
	r.Run(":8080")
}
