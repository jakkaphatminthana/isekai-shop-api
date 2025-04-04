package repository

import (
	"github.com/jakkaphatminthana/isekai-shop-api/databases"
	"github.com/jakkaphatminthana/isekai-shop-api/entities"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	_itemShopException "github.com/jakkaphatminthana/isekai-shop-api/pkg/itemShop/exception"
	_itemShopModel "github.com/jakkaphatminthana/isekai-shop-api/pkg/itemShop/model"
)

type itemShopRepositoryImpl struct {
	db     databases.Database
	logger echo.Logger
}

func NewItemShopRepositoryImpl(db databases.Database, logger echo.Logger) ItemShopRepository {
	return &itemShopRepositoryImpl{db, logger}
}

// implement
func (r *itemShopRepositoryImpl) TransactionBegin() *gorm.DB {
	tx := r.db.Connect()
	return tx.Begin()
}

// implement
func (r *itemShopRepositoryImpl) TransactionRollback(tx *gorm.DB) error {
	return tx.Rollback().Error
}

// implement
func (r *itemShopRepositoryImpl) TransactionCommit(tx *gorm.DB) error {
	return tx.Commit().Error
}

// implement
func (r *itemShopRepositoryImpl) Listing(itemFilter *_itemShopModel.ItemFilter) ([]*entities.Item, error) {
	itemList := make([]*entities.Item, 0)

	query := r.db.Connect().Model(&entities.Item{}) //select * from items

	// 1. filter data
	// 1.1 isDisable = false
	query = query.Where("is_archive = ?", false)

	// 1.2 query field custom
	if itemFilter.Name != "" {
		query = query.Where("name ilike ?", "%"+itemFilter.Name+"%")
	}
	if itemFilter.Description != "" {
		query = query.Where("description ilike ?", "%"+itemFilter.Description+"%")
	}

	// 1.3 pagination
	offset := int((itemFilter.Page - 1) * itemFilter.Size)
	limit := int(itemFilter.Size)

	// 2. find database
	// - filter field custom
	// - pagination (offset, limit)
	// - sort asc by id
	if err := query.Offset(offset).Limit(limit).Find(&itemList).Order("id asc").Error; err != nil {
		r.logger.Errorf("Failed to list items: %s", err.Error())
		return nil, &_itemShopException.ItemListing{}
	}

	return itemList, nil
}

// implement
func (r *itemShopRepositoryImpl) Counting(itemFilter *_itemShopModel.ItemFilter) (int64, error) {
	query := r.db.Connect().Model(&entities.Item{}) //select * from items

	// 1. filter data
	// 1.1 isDisable = false
	query = query.Where("is_archive = ?", false)

	// 1.2 query field custom
	if itemFilter.Name != "" {
		query = query.Where("name ilike ?", "%"+itemFilter.Name+"%")
	}
	if itemFilter.Description != "" {
		query = query.Where("description ilike ?", "%"+itemFilter.Description+"%")
	}

	// 2. counting
	var count int64

	if err := query.Count(&count).Error; err != nil {
		r.logger.Errorf("Counting items failed: %s", err.Error())
		return -1, &_itemShopException.ItemCounting{}
	}

	return count, nil
}

// implement
func (r *itemShopRepositoryImpl) FindByID(itemID uint64) (*entities.Item, error) {
	item := new(entities.Item)

	if err := r.db.Connect().First(item, itemID).Error; err != nil {
		r.logger.Errorf("Failed to find item by ID: %s", err.Error())
		return nil, &_itemShopException.ItemNotFound{}
	}

	return item, nil
}

// implement
func (r *itemShopRepositoryImpl) FindByIDList(itemIDs []uint64) ([]*entities.Item, error) {
	items := make([]*entities.Item, 0)

	if err := r.db.Connect().Model(
		&entities.Item{},
	).Where(
		"id in ?", itemIDs,
	).Find(&items).Error; err != nil {
		r.logger.Errorf("Failed to find items by ID list: %s", err.Error())
		return nil, &_itemShopException.ItemListing{}
	}

	return items, nil
}

// implementation
func (r *itemShopRepositoryImpl) PurchaseHistoryRecording(tx *gorm.DB, purchasingEntity *entities.PurchaseHistory) (*entities.PurchaseHistory, error) {
	// transaction
	conn := r.db.Connect()
	if tx != nil {
		conn = tx
	}

	insertedPurchasing := new(entities.PurchaseHistory)

	if err := conn.Create(purchasingEntity).Scan(insertedPurchasing).Error; err != nil {
		r.logger.Errorf("Failed to record purchase history: %s", err.Error())
		return nil, &_itemShopException.HistoryOfPurchaseRecording{}
	}

	return insertedPurchasing, nil
}
