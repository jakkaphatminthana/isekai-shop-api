package service

import (
	"github.com/jakkaphatminthana/isekai-shop-api/entities"
	_inventoryModel "github.com/jakkaphatminthana/isekai-shop-api/pkg/inventory/model"
	_inventoryRepository "github.com/jakkaphatminthana/isekai-shop-api/pkg/inventory/repository"
	_itemShopRepository "github.com/jakkaphatminthana/isekai-shop-api/pkg/itemShop/repository"
)

type inventoryServiceImpl struct {
	inventoryRepository _inventoryRepository.InventoryRepository
	itemShopRepository  _itemShopRepository.ItemShopRepository
}

func NewInventoryService(
	inventoryRepository _inventoryRepository.InventoryRepository,
	itemShopRepository _itemShopRepository.ItemShopRepository,
) InventoryService {
	return &inventoryServiceImpl{
		inventoryRepository: inventoryRepository,
		itemShopRepository:  itemShopRepository,
	}
}

// implementation
func (s *inventoryServiceImpl) Listing(playerID string) ([]*_inventoryModel.Inventory, error) {
	inventoryEntities, err := s.inventoryRepository.Listing(playerID)
	if err != nil {
		return nil, err
	}

	uniqueItemWithQuantityCounterList := s.getUniqueItemWithQuantityCounterList(inventoryEntities)

	return s.buildInventoryListingResult(uniqueItemWithQuantityCounterList), nil
}

// counting itemA x3 (example : data = [itemA, itemA, itemA, itemB, ...])
func (s *inventoryServiceImpl) getUniqueItemWithQuantityCounterList(
	inventoryEntities []*entities.Inventory,
) []_inventoryModel.ItemQuantityCouning {
	itemQuantityCounterList := make([]_inventoryModel.ItemQuantityCouning, 0)
	itemMapWithQuantity := make(map[uint64]uint)

	// counting (itemA: itemA++)
	for _, inventory := range inventoryEntities {
		itemMapWithQuantity[inventory.ItemID]++
	}

	// convert map to slice
	for itemID, quantity := range itemMapWithQuantity {
		itemQuantityCounterList = append(itemQuantityCounterList, _inventoryModel.ItemQuantityCouning{
			ItemID:   itemID,
			Quantity: quantity,
		})
	}

	return itemQuantityCounterList
}

func (s *inventoryServiceImpl) buildInventoryListingResult(
	uniqueItemWithQuantityCountingList []_inventoryModel.ItemQuantityCouning,
) []*_inventoryModel.Inventory {

	//get item id unique
	uniqueItemIDList := s.getItemID(uniqueItemWithQuantityCountingList)

	// []itemIDs -> []itemModels
	itemEntities, err := s.itemShopRepository.FindByIDList(uniqueItemIDList)
	if err != nil {
		return make([]*_inventoryModel.Inventory, 0)
	}

	// Map itemId with quantity (ex. {itemA: 3})
	itemMapWithQuantity := s.getItemMapWithQuantity(uniqueItemWithQuantityCountingList)

	results := make([]*_inventoryModel.Inventory, 0)
	for _, itemEntity := range itemEntities {
		results = append(results, &_inventoryModel.Inventory{
			Item:     itemEntity.ToItemModel(),
			Quantity: itemMapWithQuantity[itemEntity.ID],
		})
	}
	return results
}

// get list itemIDs
// example: [itemA, itemB,..., itemG]
func (s *inventoryServiceImpl) getItemID(
	uniqueItemWithQuantityCountingList []_inventoryModel.ItemQuantityCouning,
) []uint64 {
	uniqueItemIDList := make([]uint64, 0)

	for _, inventory := range uniqueItemWithQuantityCountingList {
		uniqueItemIDList = append(uniqueItemIDList, inventory.ItemID)
	}

	return uniqueItemIDList
}

// get Map item
// example: {potionMana: 2}
func (s *inventoryServiceImpl) getItemMapWithQuantity(
	uniqueItemWithQuantityCountingList []_inventoryModel.ItemQuantityCouning,
) map[uint64]uint {
	itemMapWithQuantity := make(map[uint64]uint)

	for _, inventory := range uniqueItemWithQuantityCountingList {
		itemMapWithQuantity[inventory.ItemID] = inventory.Quantity
	}

	return itemMapWithQuantity
}
