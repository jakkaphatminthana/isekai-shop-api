package service

import (
	"github.com/jakkaphatminthana/isekai-shop-api/entities"
	__inventoryRepository "github.com/jakkaphatminthana/isekai-shop-api/pkg/inventory/repository"
	_itemShopException "github.com/jakkaphatminthana/isekai-shop-api/pkg/itemShop/exception"
	_itemShopModel "github.com/jakkaphatminthana/isekai-shop-api/pkg/itemShop/model"
	_itemShopRepository "github.com/jakkaphatminthana/isekai-shop-api/pkg/itemShop/repository"
	_playerCoinModel "github.com/jakkaphatminthana/isekai-shop-api/pkg/playerCoin/model"
	_playerCoinRepository "github.com/jakkaphatminthana/isekai-shop-api/pkg/playerCoin/repository"
	"github.com/labstack/echo/v4"
)

type itemShopServiceImpl struct {
	itemShopRepository   _itemShopRepository.ItemShopRepository
	playerCoinRepository _playerCoinRepository.PlayerCoinRepository
	inventoryRepository  __inventoryRepository.InventoryRepository
	logger               echo.Logger
}

func NewItemShopServiceImpl(
	itemShopRepository _itemShopRepository.ItemShopRepository,
	playerCoinRepository _playerCoinRepository.PlayerCoinRepository,
	inventoryRepository __inventoryRepository.InventoryRepository,
	logger echo.Logger,
) ItemShopService {
	return &itemShopServiceImpl{
		itemShopRepository,
		playerCoinRepository,
		inventoryRepository,
		logger,
	}
}

// implement
func (s *itemShopServiceImpl) Listing(itemFilter *_itemShopModel.ItemFilter) (*_itemShopModel.ItemResult, error) {
	// Default pagination
	if itemFilter.Size <= 0 {
		itemFilter.Size = 10
	}
	if itemFilter.Page <= 0 {
		itemFilter.Page = 1
	}

	// Filter
	itemList, err := s.itemShopRepository.Listing(itemFilter)
	if err != nil {
		return nil, err
	}

	// Counting (total item)
	itemCounting, err := s.itemShopRepository.Counting(itemFilter)
	if err != nil {
		return nil, err
	}

	// Pagination calculate
	size := itemFilter.Size
	page := itemFilter.Page
	totalPage := s.totalPageCalculation(itemCounting, size)

	result := s.toItemResultResponse(itemList, page, totalPage)

	return result, nil
}

// implement
// 1. Find item by ID
// 2. Total price calculation
// 3. Check player coin
// 4. Purchase history recording
// 5. Coin deducting
// 6. Inventory filling
// 7. Return player coin
func (s *itemShopServiceImpl) Buying(buyingReq *_itemShopModel.BuyingReq) (*_playerCoinModel.PlayerCoin, error) {
	// 1. Find item by ID
	itemEntity, err := s.itemShopRepository.FindByID(buyingReq.ItemID)
	if err != nil {
		return nil, err
	}

	// 2. Total price calculation
	totalPrice := s.totalPriceCalculation(itemEntity.ToItemModel(), buyingReq.Quantity)

	// 3. Check player coin
	if err := s.playerCoinChecking(buyingReq.PlayerID, totalPrice); err != nil {
		return nil, err
	}

	// transaction begin
	tx := s.itemShopRepository.TransactionBegin()

	// 4. Purchase history recording
	purchaseRecording, err := s.itemShopRepository.PurchaseHistoryRecording(tx, &entities.PurchaseHistory{
		PlayerID:        buyingReq.PlayerID,
		ItemID:          buyingReq.ItemID,
		ItemName:        itemEntity.Name,
		ItemDescription: itemEntity.Description,
		ItemPrice:       itemEntity.Price,
		ItemPicture:     itemEntity.Picture,
		Quantity:        buyingReq.Quantity,
		IsBuying:        true,
	})
	if err != nil {
		s.itemShopRepository.TransactionRollback(tx)
		return nil, err
	}
	s.logger.Infof("Purchase history recored: %s", purchaseRecording.ID)

	// 5. Coin deducting
	playerCoin, err := s.playerCoinRepository.CoinAdding(tx, &entities.PlayerCoin{
		PlayerID: buyingReq.PlayerID,
		Amount:   -totalPrice,
	})
	if err != nil {
		s.itemShopRepository.TransactionRollback(tx)
		return nil, err
	}
	s.logger.Infof("Player coin deducted: %d", playerCoin.Amount)

	// 6. Inventory filling
	inventoryEntity, err := s.inventoryRepository.Filling(
		tx,
		buyingReq.PlayerID,
		buyingReq.ItemID,
		int(buyingReq.Quantity),
	)
	if err != nil {
		s.itemShopRepository.TransactionRollback(tx)
		return nil, err
	}
	s.logger.Infof("Inventory filled: %d", len(inventoryEntity))

	// commit transaction
	if err := s.itemShopRepository.TransactionCommit(tx); err != nil {
		return nil, err
	}

	// 7. Return player coin
	return playerCoin.ToPlayerCoinModel(), nil
}

// implement
// 1. Find item by ID
// 2. Total price calculation
// 3. Check player item
// 4. Purchase history recording
// 5. Coin adding
// 6. Inventory removing
// 7. Return player coin
func (s *itemShopServiceImpl) Selling(sellingReq *_itemShopModel.SellingReq) (*_playerCoinModel.PlayerCoin, error) {
	// 1. Find item by ID
	itemEntity, err := s.itemShopRepository.FindByID(sellingReq.ItemID)
	if err != nil {
		return nil, err
	}

	// 2. Total price calculation
	totalPrice := s.totalPriceCalculation(itemEntity.ToItemModel(), sellingReq.Quantity)
	totalPrice = totalPrice / 2 //fix price XD

	// 3. Check player item
	if err := s.playerItemChecking(sellingReq.PlayerID, sellingReq.ItemID, sellingReq.Quantity); err != nil {
		return nil, err
	}

	// transaction begin
	tx := s.itemShopRepository.TransactionBegin()

	// 4. Purchase history recording
	purchaseRecording, err := s.itemShopRepository.PurchaseHistoryRecording(tx, &entities.PurchaseHistory{
		PlayerID:        sellingReq.PlayerID,
		ItemID:          sellingReq.ItemID,
		ItemName:        itemEntity.Name,
		ItemDescription: itemEntity.Description,
		ItemPrice:       itemEntity.Price,
		ItemPicture:     itemEntity.Picture,
		Quantity:        sellingReq.Quantity,
		IsBuying:        false,
	})
	if err != nil {
		s.itemShopRepository.TransactionRollback(tx)
		return nil, err
	}
	s.logger.Infof("Purchase history recored: %s", purchaseRecording.ID)

	// 5. Coin adding
	playerCoin, err := s.playerCoinRepository.CoinAdding(tx, &entities.PlayerCoin{
		PlayerID: sellingReq.PlayerID,
		Amount:   totalPrice,
	})
	if err != nil {
		s.itemShopRepository.TransactionRollback(tx)
		return nil, err
	}
	s.logger.Infof("Player coin adding: %d", playerCoin.Amount)

	// 6. Inventory removing
	if err := s.inventoryRepository.Removing(
		tx,
		sellingReq.PlayerID,
		sellingReq.ItemID,
		int(sellingReq.Quantity),
	); err != nil {
		s.itemShopRepository.TransactionRollback(tx)
		return nil, err
	}
	s.logger.Infof("Inventory itemID: %d,removing: %d", sellingReq.ItemID, sellingReq.Quantity)

	// transaction commit
	if err := s.itemShopRepository.TransactionCommit(tx); err != nil {
		return nil, err
	}

	// 7. Return player coin
	return playerCoin.ToPlayerCoinModel(), nil
}

func (s *itemShopServiceImpl) totalPageCalculation(totalItems int64, size int64) int64 {
	totalPage := totalItems / size

	// case: 11/5 = 2.1 -> 3 (+ 1 page)
	if totalItems%size != 0 {
		totalPage++
	}

	return totalPage
}

func (s *itemShopServiceImpl) toItemResultResponse(itemEntityList []*entities.Item, page, totalPage int64) *_itemShopModel.ItemResult {
	// Mapper entity to model
	itemModelList := make([]*_itemShopModel.Item, 0)
	for _, item := range itemEntityList {
		itemModelList = append(itemModelList, item.ToItemModel())
	}

	return &_itemShopModel.ItemResult{
		Items: itemModelList,
		Paginate: _itemShopModel.PaginateResult{
			Page:      page,
			TotalPage: totalPage,
		},
	}
}

func (s *itemShopServiceImpl) totalPriceCalculation(item *_itemShopModel.Item, qty uint) int64 {
	return int64(item.Price) * int64(qty)
}

func (s *itemShopServiceImpl) playerCoinChecking(playerID string, totalPrice int64) error {
	playerCoin, err := s.playerCoinRepository.Showing(playerID)
	if err != nil {
		return err
	}

	if playerCoin.Coin < totalPrice {
		s.logger.Error("Player coin is not enough")
		return &_itemShopException.CoinNotEnough{}
	}

	return nil
}

func (s *itemShopServiceImpl) playerItemChecking(playerID string, itemID uint64, qty uint) error {
	itemCounting := s.inventoryRepository.PlayerItemCounting(playerID, itemID)

	if int(itemCounting) < int(qty) {
		s.logger.Error("Player item is not enough")
		return &_itemShopException.ItemNotEnough{ItemID: itemID}
	}

	return nil
}
