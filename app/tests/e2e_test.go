// app/e2e_test.go
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"stock-simulator/internal/handlers"
	"stock-simulator/internal/services"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

// setupRouter configuration Miniredis + Gin + Handlery
func setupRouter() (*gin.Engine, *miniredis.Miniredis) {
	gin.SetMode(gin.TestMode)

	// create redis db in ram
	mr, _ := miniredis.Run()
	rdb := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	walletService := services.NewWalletService(rdb)
	walletHandler := handlers.NewWalletHandler(walletService)

	r := gin.Default()
	r.POST("/stocks", walletHandler.SetBankState)
	r.POST("/wallets/:wallet_id/stocks/:stock_name", walletHandler.HandleStockOperation)
	r.GET("/wallets/:wallet_id", walletHandler.GetWallet)

	return r, mr
}

func TestStockMarketFlow(t *testing.T) {
	r, mr := setupRouter()
	defer mr.Close()

	// 1 - add stocks to back
	t.Run("Set Bank State", func(t *testing.T) {
		body := `{"stocks": [{"name":"AAPL", "quantity":1}]}`
		req, _ := http.NewRequest("POST", "/stocks", bytes.NewBufferString(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// 2 - buy stocks from back
	t.Run("Buy Stock - Success", func(t *testing.T) {
		body := `{"type":"buy"}`
		req, _ := http.NewRequest("POST", "/wallets/w1/stocks/AAPL", bytes.NewBufferString(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// 3 - buy stock without money - 400
	t.Run("Buy Stock - Insufficient Funds (400)", func(t *testing.T) {
		body := `{"type":"buy"}`
		req, _ := http.NewRequest("POST", "/wallets/w2/stocks/AAPL", bytes.NewBufferString(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// 4 - unknown stock - 404
	t.Run("Buy Stock - Not Found (404)", func(t *testing.T) {
		body := `{"type":"buy"}`
		req, _ := http.NewRequest("POST", "/wallets/w1/stocks/MSFT", bytes.NewBufferString(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	// sell proper stock - 200
	t.Run("Sell Stock - Success", func(t *testing.T) {
		body := `{"type":"sell"}`
		req, _ := http.NewRequest("POST", "/wallets/w1/stocks/AAPL", bytes.NewBufferString(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// 6 - sell stock with empty wallet - 400
	t.Run("Sell Stock - Empty Wallet (400)", func(t *testing.T) {
		body := `{"type":"sell"}`
		req, _ := http.NewRequest("POST", "/wallets/w1/stocks/AAPL", bytes.NewBufferString(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// 7 - check wallet status (should be empty here)
	t.Run("Check Wallet", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/wallets/w1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)

		// Wymagane sprawdzanie tablicy akcji wg wytycznych z zadania
		stocks := resp["stocks"].([]interface{})
		assert.Equal(t, 0, len(stocks))
	})
}
