package repository

import (
	"github.com/jakkaphatminthana/isekai-shop-api/databases"
	"github.com/jakkaphatminthana/isekai-shop-api/entities"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	_inventotyException "github.com/jakkaphatminthana/isekai-shop-api/pkg/inventory/exception"
)

type inventoryRepositoryImpl struct {
	db     databases.Database
	logger echo.Logger
}

func NewInventoryRepository(db databases.Database, logger echo.Logger) InventoryRepository {
	return &inventoryRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

// implementation
func (r *inventoryRepositoryImpl) Filling(tx *gorm.DB, playerID string, itemID uint64, qty int) ([]*entities.Inventory, error) {
	conn := r.db.Connect()
	if tx != nil {
		conn = tx
	}

	inventoryEntities := make([]*entities.Inventory, 0)

	// loop inventory by qty
	for range qty {
		inventoryEntities = append(inventoryEntities, &entities.Inventory{
			PlayerID: playerID,
			ItemID:   itemID,
		})
	}

	if err := conn.Create(inventoryEntities).Error; err != nil {
		r.logger.Errorf("error  filling inventory: %s", err.Error())
		return nil, &_inventotyException.InventoryFilling{
			PlayerID: playerID,
			ItemID:   itemID,
		}
	}

	return inventoryEntities, nil
}

// implementation
func (r *inventoryRepositoryImpl) Removing(tx *gorm.DB, playerID string, itemID uint64, limit int) error {
	conn := r.db.Connect()
	if tx != nil {
		conn = tx
	}

	inventoryEntities, err := r.findItemInInventoryByID(playerID, itemID, limit)
	if err != nil {
		return err
	}

	for _, inventory := range inventoryEntities {
		inventory.IsDeleted = true

		if err := conn.Model(
			&entities.Inventory{},
		).Where(
			"id = ?", inventory.ID,
		).Updates(
			inventory,
		).Error; err != nil {
			tx.Rollback()
			r.logger.Errorf("error removing player item in inventory: %s", err.Error())
			return &_inventotyException.PlayerItemRemoving{
				ItemID: itemID,
			}
		}
	}

	return nil
}

// implementation
func (r *inventoryRepositoryImpl) PlayerItemCounting(playerID string, itemID uint64) int64 {
	var count int64

	if err := r.db.Connect().Model(
		&entities.Inventory{},
	).Where(
		"player_id = ? AND item_id = ? AND is_deleted = ?", playerID, itemID, false,
	).Count(&count).Error; err != nil {
		r.logger.Errorf("error counting player item in inventory: %s", err.Error())
		return -1
	}
	return count
}

// implementation
func (r *inventoryRepositoryImpl) Listing(playerID string) ([]*entities.Inventory, error) {
	inventoryEntities := make([]*entities.Inventory, 0)

	if err := r.db.Connect().Where(
		"player_id = ? AND is_deleted = ?", playerID, false,
	).Find(&inventoryEntities).Error; err != nil {
		r.logger.Errorf("error listing player inventory: %s", err.Error())
		return nil, &_inventotyException.PlayerItemFinding{
			PlayerID: playerID,
		}
	}
	return inventoryEntities, nil
}

func (r *inventoryRepositoryImpl) findItemInInventoryByID(
	playerID string,
	itemID uint64,
	limit int,
) ([]*entities.Inventory, error) {
	inventoryEntities := make([]*entities.Inventory, 0)

	if err := r.db.Connect().Where(
		"player_id = ? AND item_id = ? AND is_deleted = ?", playerID, itemID, false,
	).Limit(
		limit,
	).Find(&inventoryEntities).Error; err != nil {
		r.logger.Errorf("error finding player item in inventory by ID: %s", err.Error())
		return nil, &_inventotyException.PlayerItemRemoving{
			ItemID: itemID,
		}
	}

	return inventoryEntities, nil
}
