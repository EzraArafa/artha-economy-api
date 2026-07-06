package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/EzraArafa/artha-economy-api/model"
	"github.com/EzraArafa/artha-economy-api/service"
)

type ItemController struct {
	itemService *service.ItemService
}

func NewItemController(itemService *service.ItemService) *ItemController {
	return &ItemController{itemService: itemService}
}

func (ctrl *ItemController) Create(c *gin.Context) {
	var item model.Item

	//1. Menagkap data JSON (misal: nama dan harga barang)
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak sesuai"})
		return
	}

	//2. Lempar ke service untuk di validasi
	err := ctrl.itemService.CreateItem(&item)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//3. Kembalikan respon sukses
	c.JSON(http.StatusCreated, gin.H{
		"message": "Barang berhasil ditambahkan ke database!",
		"data":    item,
	})
}
