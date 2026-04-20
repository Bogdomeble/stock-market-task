package main

import (

	"stock-simulator/internal/handlers"
    "stock-simulator/internal/repository"
    "stock-simulator/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// init redis
	rdb := repository.NewRedisClient()

	// init services
	walletService := services.NewWalletService(rdb)

	// init handlers
	walletHandler := handlers.NewWalletHandler(walletService)

	r := gin.Default()

	// routes
	r.GET("/health", handlers.HealthCheck)
	r.POST("/chaos", handlers.Chaos)

	r.POST("/wallets/:wallet_id/stocks/:stock_name", walletHandler.HandleStockOperation)
	r.GET("/wallets/:wallet_id", walletHandler.GetWallet)

	r.Run(":8080")
}