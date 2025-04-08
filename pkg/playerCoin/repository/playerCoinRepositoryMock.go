package repository

import (
	"github.com/jakkaphatminthana/isekai-shop-api/entities"
	_playerCoinModel "github.com/jakkaphatminthana/isekai-shop-api/pkg/playerCoin/model"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type PlayerCoinRepositoryMock struct {
	mock.Mock
}

func (m *PlayerCoinRepositoryMock) CoinAdding(tx *gorm.DB, playerCoinEntity *entities.PlayerCoin) (*entities.PlayerCoin, error) {
	// 0 = return index 0 : *entities.PlayerCoin
	// 1 = return index 1 : error
	args := m.Called(tx, playerCoinEntity)
	return args.Get(0).(*entities.PlayerCoin), args.Error(1)
}

func (m *PlayerCoinRepositoryMock) Showing(playerID string) (*_playerCoinModel.PlayerCoinShowing, error) {
	// 0 = return index 0 : *_playerCoinModel.PlayerCoinShowing
	// 1 = return index 1 : error
	args := m.Called(playerID)
	return args.Get(0).(*_playerCoinModel.PlayerCoinShowing), args.Error(1)
}
