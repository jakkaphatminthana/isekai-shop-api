package entities

import (
	"time"

	_itemShopModel "github.com/jakkaphatminthana/isekai-shop-api/pkg/itemShop/model"
)

type (
	Item struct {
		ID          uint64    `gorm:"primaryKey;autoIncrement;"`
		AdminID     *string   `gorm:"type:varchar(64);"`
		Name        string    `gorm:"type:varchar(64);unique;not null;"`
		Description string    `gorm:"type:varchar(128);not null;"`
		Picture     string    `gorm:"type:varchar(256);not null;"`
		Price       uint      `gorm:"not null;"`
		IsArchive   bool      `gorm:"not null;default:false;"`
		CreatedAt   time.Time `gorm:"autoCreateTime;not null;"`
		UpdatedAt   time.Time `gorm:"autoUpdateTime;not null;"`
	}
)

func (i *Item) ToItemModel() *_itemShopModel.Item {
	return &_itemShopModel.Item{
		ID:          i.ID,
		Name:        i.Name,
		Description: i.Description,
		Picture:     i.Picture,
		Price:       i.Price,
	}
}
