package model

import (
	_itemShopModel "github.com/jakkaphatminthana/isekai-shop-api/pkg/itemShop/model"
)

type (
	Inventory struct {
		Item     *_itemShopModel.Item `json:"item"`
		Quantity uint                 `json:"quantity"`
	}

	ItemQuantityCouning struct {
		ItemID   uint64
		Quantity uint
	}
)
