package handlers

import (
	"net/http"

	"stock-simulator/internal/services"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	service *services.WalletService
}

func NewWalletHandler(s *services.WalletService) *WalletHandler {
	return &WalletHandler{service: s}
}

func (h *WalletHandler) HandleStockOperation(c *gin.Context) {
	walletID := c.Param("wallet_id")
	stockName := c.Param("stock_name")

	var body struct {
		Type string `json:"type"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	var err error

	if body.Type == "buy" {
		err = h.service.Buy(walletID, stockName)
	} else if body.Type == "sell" {
		err = h.service.Sell(walletID, stockName)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid type"})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *WalletHandler) GetWallet(c *gin.Context) {
	// TODO: implement later
	c.JSON(200, gin.H{"message": "not implemented yet"})
}