package repository

import (
	"github.com/jakkaphatminthana/isekai-shop-api/entities"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type InventoryRepositoryMock struct {
	mock.Mock
}

func (m *InventoryRepositoryMock) Filling(tx *gorm.DB, playerID string, itemID uint64, qty int) ([]*entities.Inventory, error) {
	// 0 = return index 0 : []*entities.Inventory
	// 1 = return index 1 : error
	args := m.Called(tx, playerID, itemID, qty)
	return args.Get(0).([]*entities.Inventory), args.Error(1)
}

func (m *InventoryRepositoryMock) Removing(tx *gorm.DB, playerID string, itemID uint64, limit int) error {
	// 0 = return index 0 : error
	args := m.Called(tx, playerID, itemID, limit)
	return args.Error(0)
}

func (m *InventoryRepositoryMock) PlayerItemCounting(playerID string, itemID uint64) int64 {
	// 0 = return index 0 : int64
	args := m.Called(playerID, itemID)
	return args.Get(0).(int64)
}

func (m *InventoryRepositoryMock) Listing(playerID string) ([]*entities.Inventory, error) {
	// 0 = return index 0 : []*entities.Inventory
	// 1 = return index 1 : error
	args := m.Called(playerID)
	return args.Get(0).([]*entities.Inventory), args.Error(1)
}
