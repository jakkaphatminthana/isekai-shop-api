package tests

import (
	"testing"

	"github.com/jakkaphatminthana/isekai-shop-api/entities"
	_inventoryRepository "github.com/jakkaphatminthana/isekai-shop-api/pkg/inventory/repository"
	_itemShopException "github.com/jakkaphatminthana/isekai-shop-api/pkg/itemShop/exception"
	_itemShopModel "github.com/jakkaphatminthana/isekai-shop-api/pkg/itemShop/model"
	_itemShopRepository "github.com/jakkaphatminthana/isekai-shop-api/pkg/itemShop/repository"
	_itemShopService "github.com/jakkaphatminthana/isekai-shop-api/pkg/itemShop/service"
	_playerCoinModel "github.com/jakkaphatminthana/isekai-shop-api/pkg/playerCoin/model"
	_playerCoinRepository "github.com/jakkaphatminthana/isekai-shop-api/pkg/playerCoin/repository"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestItemBuyingSuccess(t *testing.T) {
	itemShopRepositoryMock := new(_itemShopRepository.ItemShopRepositoryMock)
	playerCoinRepositoryMock := new(_playerCoinRepository.PlayerCoinRepositoryMock)
	inventoryRepositoryMock := new(_inventoryRepository.InventoryRepositoryMock)
	echoLogger := echo.New().Logger

	itemShopService := _itemShopService.NewItemShopServiceImpl(
		itemShopRepositoryMock,
		playerCoinRepositoryMock,
		inventoryRepositoryMock,
		echoLogger,
	)

	tx := &gorm.DB{}
	itemShopRepositoryMock.On("TransactionBegin").Return(tx)
	itemShopRepositoryMock.On("TransactionCommit", tx).Return(nil)
	itemShopRepositoryMock.On("TransactionRollback", tx).Return(nil)

	itemShopRepositoryMock.On("FindByID", uint64(1)).Return(&entities.Item{
		ID:          1,
		Name:        "Sword of Tester",
		Price:       1000,
		Description: "A sword that can be used to test the enemy's defense",
		Picture:     "https://www.google.com/sword-of-tester.jpg",
	}, nil)

	playerCoinRepositoryMock.On("Showing", "P001").Return(&_playerCoinModel.PlayerCoinShowing{
		PlayerID: "P001",
		Coin:     5000,
	}, nil)

	itemShopRepositoryMock.On("PurchaseHistoryRecording", tx, &entities.PurchaseHistory{
		PlayerID:        "P001",
		ItemID:          1,
		ItemName:        "Sword of Tester",
		ItemDescription: "A sword that can be used to test the enemy's defense",
		ItemPicture:     "https://www.google.com/sword-of-tester.jpg",
		ItemPrice:       1000,
		Quantity:        3,
		IsBuying:        true,
	}).Return(&entities.PurchaseHistory{
		PlayerID:        "P001",
		ItemID:          1,
		ItemName:        "Sword of Tester",
		ItemDescription: "A sword that can be used to test the enemy's defense",
		ItemPicture:     "https://www.google.com/sword-of-tester.jpg",
		ItemPrice:       1000,
		Quantity:        3,
		IsBuying:        true,
	}, nil)

	inventoryRepositoryMock.On("Filling", tx, "P001", uint64(1), int(3)).Return([]*entities.Inventory{
		{
			PlayerID: "P001",
			ItemID:   1,
		},
		{
			PlayerID: "P001",
			ItemID:   1,
		},
		{
			PlayerID: "P001",
			ItemID:   1,
		},
	}, nil)

	playerCoinRepositoryMock.On("CoinAdding", tx, &entities.PlayerCoin{
		PlayerID: "P001",
		Amount:   -3000,
	}).Return(&entities.PlayerCoin{
		PlayerID: "P001",
		Amount:   -3000,
	}, nil)

	// Start of test case
	type args struct {
		label    string
		in       *_itemShopModel.BuyingReq
		expected *_playerCoinModel.PlayerCoin
	}

	cases := []args{
		{
			label: "Test case 1 | Buying item success",
			in: &_itemShopModel.BuyingReq{
				PlayerID: "P001",
				ItemID:   1,
				Quantity: 3,
			},
			expected: &_playerCoinModel.PlayerCoin{
				PlayerID: "P001",
				Amount:   -3000,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			result, err := itemShopService.Buying(c.in)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.EqualValues(t, c.expected, result)
		})
	}
}

func TestItemBuyingFail(t *testing.T) {
	itemShopRepositoryMock := new(_itemShopRepository.ItemShopRepositoryMock)
	playerCoinRepositoryMock := new(_playerCoinRepository.PlayerCoinRepositoryMock)
	inventoryRepositoryMock := new(_inventoryRepository.InventoryRepositoryMock)
	echoLogger := echo.New().Logger

	itemShopService := _itemShopService.NewItemShopServiceImpl(
		itemShopRepositoryMock,
		playerCoinRepositoryMock,
		inventoryRepositoryMock,
		echoLogger,
	)

	tx := &gorm.DB{}
	itemShopRepositoryMock.On("TransactionBegin").Return(tx)
	itemShopRepositoryMock.On("TransactionCommit", tx).Return(nil)
	itemShopRepositoryMock.On("TransactionRollback", tx).Return(nil)

	itemShopRepositoryMock.On("FindByID", uint64(1)).Return(&entities.Item{
		ID:          1,
		Name:        "Sword of Tester",
		Price:       1000,
		Description: "A sword that can be used to test the enemy's defense",
		Picture:     "https://www.google.com/sword-of-tester.jpg",
	}, nil)

	playerCoinRepositoryMock.On("Showing", "P001").Return(&_playerCoinModel.PlayerCoinShowing{
		PlayerID: "P001",
		Coin:     1000,
	}, nil)

	itemShopRepositoryMock.On("PurchaseHistoryRecording", tx, &entities.PurchaseHistory{
		PlayerID:        "P001",
		ItemID:          1,
		ItemName:        "Sword of Tester",
		ItemDescription: "A sword that can be used to test the enemy's defense",
		ItemPicture:     "https://www.google.com/sword-of-tester.jpg",
		ItemPrice:       1000,
		Quantity:        3,
		IsBuying:        true,
	}).Return(&entities.PurchaseHistory{
		PlayerID:        "P001",
		ItemID:          1,
		ItemName:        "Sword of Tester",
		ItemDescription: "A sword that can be used to test the enemy's defense",
		ItemPicture:     "https://www.google.com/sword-of-tester.jpg",
		ItemPrice:       1000,
		Quantity:        3,
		IsBuying:        true,
	}, nil)

	inventoryRepositoryMock.On("Filling", tx, "P001", uint64(1), int(3)).Return([]*entities.Inventory{
		{
			PlayerID: "P001",
			ItemID:   1,
		},
		{
			PlayerID: "P001",
			ItemID:   1,
		},
		{
			PlayerID: "P001",
			ItemID:   1,
		},
	}, nil)

	playerCoinRepositoryMock.On("CoinAdding", tx, &entities.PlayerCoin{
		PlayerID: "P001",
		Amount:   -3000,
	}).Return(&entities.PlayerCoin{
		PlayerID: "P001",
		Amount:   -3000,
	}, nil)

	// Start of test case
	type args struct {
		label    string
		in       *_itemShopModel.BuyingReq
		expected error
	}

	cases := []args{
		{
			label: "Test case 2 | Buying item failed because player coin not enough",
			in: &_itemShopModel.BuyingReq{
				PlayerID: "P001",
				ItemID:   1,
				Quantity: 3,
			},
			expected: &_itemShopException.CoinNotEnough{},
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			result, err := itemShopService.Buying(c.in)
			assert.Nil(t, result)
			assert.Error(t, err)
			assert.EqualValues(t, c.expected, err)
		})
	}
}
