package service

import (
	_inventoryModel "github.com/jakkaphatminthana/isekai-shop-api/pkg/inventory/model"
)

type InventoryService interface {
	Listing(playerID string) ([]*_inventoryModel.Inventory, error)
}
