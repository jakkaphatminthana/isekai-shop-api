package entities

import "time"

type Admin struct {
	ID        string    `gorm:"primaryKey;type:varchar(64);"`
	Items     []Item    `gorm:"foreignKey:AdminID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Email     string    `gorm:"type:varchar(128);unique;not null;"`
	Name      string    `gorm:"type:varchar(128);not null;"`
	Avatar    string    `gorm:"type:varchar(256);not null;default:'';"`
	CreatedAt time.Time `gorm:"autoCreateTime;not null;"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;not null;"`
}
