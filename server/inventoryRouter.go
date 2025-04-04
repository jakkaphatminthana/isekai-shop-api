package server

import (
	_inventoryController "github.com/jakkaphatminthana/isekai-shop-api/pkg/inventory/controller"
	_inventoryRepository "github.com/jakkaphatminthana/isekai-shop-api/pkg/inventory/repository"
	_inventoryService "github.com/jakkaphatminthana/isekai-shop-api/pkg/inventory/service"
	_itemShopRepository "github.com/jakkaphatminthana/isekai-shop-api/pkg/itemShop/repository"
)

func (s *echoServer) initInventoryRouter(m *authorizingMiddleware) {
	router := s.app.Group("/v1/inventory")

	inventoryRepository := _inventoryRepository.NewInventoryRepository(s.db, s.app.Logger)
	itemShopRepository := _itemShopRepository.NewItemShopRepositoryImpl(s.db, s.app.Logger)

	inventoryService := _inventoryService.NewInventoryService(
		inventoryRepository,
		itemShopRepository,
	)
	inventoryController := _inventoryController.NewInventoryController(inventoryService, s.app.Logger)

	router.GET("", inventoryController.Listing, m.PlayerAuthorizing)
}
