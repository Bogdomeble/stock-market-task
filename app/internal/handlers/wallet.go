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

type StockDTO struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
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
		if err.Error() == "NOT_FOUND" {
			c.JSON(http.StatusNotFound, gin.H{"error": "stock not found in bank"})
			return
		}
		if err.Error() == "INSUFFICIENT_FUNDS" || err.Error() == "NO_STOCK_IN_WALLET" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient stocks"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *WalletHandler) GetWallet(c *gin.Context) {
	walletID := c.Param("wallet_id")
	stocksMap, err := h.service.GetWallet(walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	stocks := []StockDTO{} // empty slice instead of nil
	for name, qty := range stocksMap {
		stocks = append(stocks, StockDTO{Name: name, Quantity: qty})
	}

	c.JSON(http.StatusOK, gin.H{"id": walletID, "stocks": stocks})
}

func (h *WalletHandler) GetWalletStockQuantity(c *gin.Context) {
	walletID := c.Param("wallet_id")
	stockName := c.Param("stock_name")

	qty, err := h.service.GetWalletStock(walletID, stockName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// returns a single number, like: 99
	c.JSON(http.StatusOK, qty)
}

func (h *WalletHandler) GetBankState(c *gin.Context) {
	stocksMap, err := h.service.GetBankState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	stocks := []StockDTO{}
	for name, qty := range stocksMap {
		stocks = append(stocks, StockDTO{Name: name, Quantity: qty})
	}
	c.JSON(http.StatusOK, gin.H{"stocks": stocks})
}

func (h *WalletHandler) SetBankState(c *gin.Context) {
	var body struct {
		Stocks []StockDTO `json:"stocks"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	stocksMap := make(map[string]int)
	for _, s := range body.Stocks {
		stocksMap[s.Name] = s.Quantity
	}

	if err := h.service.SetBankState(stocksMap); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (h *WalletHandler) GetAuditLog(c *gin.Context) {
	logs, err := h.service.GetLog()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"log": logs})
}
