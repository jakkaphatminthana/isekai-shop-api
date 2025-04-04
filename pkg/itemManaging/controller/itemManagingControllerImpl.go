package controller

import (
	"net/http"
	"strconv"

	"github.com/jakkaphatminthana/isekai-shop-api/pkg/custom"
	_itemManagingModel "github.com/jakkaphatminthana/isekai-shop-api/pkg/itemManaging/model"
	_itemManagingService "github.com/jakkaphatminthana/isekai-shop-api/pkg/itemManaging/service"
	"github.com/jakkaphatminthana/isekai-shop-api/pkg/validation"
	"github.com/labstack/echo/v4"
)

type itemManagingControllerImpl struct {
	itemManagingService _itemManagingService.ItemManagingService
}

func NewManagingControllerImpl(
	itemManagingService _itemManagingService.ItemManagingService,
) ItemManagingController {
	return &itemManagingControllerImpl{itemManagingService}
}

// implement
func (c *itemManagingControllerImpl) Creating(pctx echo.Context) error {
	//validate admin role
	adminID, err := validation.AdminIDGetting(pctx)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	itemCreatingReq := new(_itemManagingModel.ItemCreatingReq)

	// validate
	customEchoRequest := custom.NewCustomEchoRequest(pctx)
	if err := customEchoRequest.Bind(itemCreatingReq); err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}
	itemCreatingReq.AdminId = adminID

	// creating
	item, err := c.itemManagingService.Creating(itemCreatingReq)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.JSON(http.StatusCreated, item)
}

// implement
func (c *itemManagingControllerImpl) Editing(pctx echo.Context) error {
	itemEditingReq := new(_itemManagingModel.ItemEditingReq)

	itemID, err := c.getItemID(pctx)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	//validate
	customEchoRequest := custom.NewCustomEchoRequest(pctx)
	if err := customEchoRequest.Bind(itemEditingReq); err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	//editing
	item, err := c.itemManagingService.Editing(itemID, itemEditingReq)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.JSON(http.StatusOK, item)
}

// implement
func (c *itemManagingControllerImpl) Archiving(pctx echo.Context) error {
	itemEditingReq := new(_itemManagingModel.ItemEditingReq)

	itemID, err := c.getItemID(pctx)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	//validate
	customEchoRequest := custom.NewCustomEchoRequest(pctx)
	if err := customEchoRequest.Bind(itemEditingReq); err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	if err := c.itemManagingService.Archiving(itemID); err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.NoContent(http.StatusNoContent)
}

func (c *itemManagingControllerImpl) getItemID(pctx echo.Context) (uint64, error) {
	itemID := pctx.Param("itemID")

	itemIDUint64, err := strconv.ParseUint(itemID, 10, 64)
	if err != nil {
		return 0, err
	}

	return itemIDUint64, nil
}
