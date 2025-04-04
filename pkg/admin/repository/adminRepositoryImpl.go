package repository

import (
	"github.com/jakkaphatminthana/isekai-shop-api/databases"
	"github.com/jakkaphatminthana/isekai-shop-api/entities"
	"github.com/labstack/echo/v4"

	_adminException "github.com/jakkaphatminthana/isekai-shop-api/pkg/admin/exception"
)

type adminRepositoryImpl struct {
	db     databases.Database
	logger echo.Logger
}

func NewAdminRepositoryImpl(
	db databases.Database,
	logger echo.Logger,
) AdminRepository {
	return &adminRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

// implement
func (r *adminRepositoryImpl) Creating(adminEntity *entities.Admin) (*entities.Admin, error) {
	admin := new(entities.Admin)

	if err := r.db.Connect().Create(adminEntity).Scan(admin).Error; err != nil {
		r.logger.Errorf("Creating admin failed: %s", err.Error())
		return nil, &_adminException.AdminCreating{AdminID: admin.ID}
	}

	return admin, nil
}

// implement
func (r *adminRepositoryImpl) FindByID(adminID string) (*entities.Admin, error) {
	admin := new(entities.Admin)

	if err := r.db.Connect().Where("id = ?", adminID).First(admin).Error; err != nil {
		r.logger.Errorf("Find admin by ID failed: %s", err.Error())
		return nil, &_adminException.AdminNotFound{AdminID: adminID}
	}

	return admin, nil
}
