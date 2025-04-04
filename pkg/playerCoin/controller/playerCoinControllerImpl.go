package controller

import (
	"net/http"

	"github.com/jakkaphatminthana/isekai-shop-api/pkg/custom"
	_playerCoinModel "github.com/jakkaphatminthana/isekai-shop-api/pkg/playerCoin/model"
	_playerCoinService "github.com/jakkaphatminthana/isekai-shop-api/pkg/playerCoin/service"
	"github.com/jakkaphatminthana/isekai-shop-api/pkg/validation"
	"github.com/labstack/echo/v4"
)

type playerCoinControllerImpl struct {
	playerCoinService _playerCoinService.PlayerCoinService
}

func NewPlayerCoinControllerImpl(playerCoinService _playerCoinService.PlayerCoinService) PlayerCoinController {
	return &playerCoinControllerImpl{
		playerCoinService: playerCoinService,
	}
}

// implentation
func (c *playerCoinControllerImpl) CoinAdding(pctx echo.Context) error {
	// validate player
	playerID, err := validation.PlayerIDGetting(pctx)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	coinAddingReq := new(_playerCoinModel.CoinAddingReq)

	customEchoRequest := custom.NewCustomEchoRequest(pctx)
	if err := customEchoRequest.Bind(coinAddingReq); err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}
	coinAddingReq.PlayerID = playerID

	// creating
	playerCoin, err := c.playerCoinService.CoinAdding(coinAddingReq)
	if err != nil {
		return custom.Error(pctx, http.StatusInternalServerError, err)
	}

	return pctx.JSON(http.StatusCreated, playerCoin)
}

// implentation
func (c *playerCoinControllerImpl) Showing(pctx echo.Context) error {
	playerID, err := validation.PlayerIDGetting(pctx)
	if err != nil {
		return custom.Error(pctx, http.StatusBadRequest, err)
	}

	playerCoinShowing := c.playerCoinService.Showing(playerID)

	return pctx.JSON(http.StatusOK, playerCoinShowing)
}
