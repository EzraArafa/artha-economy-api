package service

import (
	"errors"
	"strings"

	"github.com/EzraArafa/artha-economy-api/model"
	"github.com/EzraArafa/artha-economy-api/repository"
)

type ItemService struct {
	itemRepo *repository.ItemRepository
}

func NewItemService(itemRepo *repository.ItemRepository) *ItemService {
	return &ItemService{itemRepo: itemRepo}
}

func (s *ItemService) CreateItem(item *model.Item) error {
	if strings.TrimSpace(item.Name) == "" {
		return errors.New("nama barang tidak boleh kosong")
	}

	if item.Price < 0 {
		return errors.New("harga barang tidak boleh minus")
	}

	return s.itemRepo.Create(item)
}

func (s *ItemService) GetItemByID(id int) (*model.Item, error) {
	return s.itemRepo.FindByID(id)
}
