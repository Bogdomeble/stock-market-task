package main

import (
	"stock-simulator/internal/handlers"
	"stock-simulator/internal/repository"
	"stock-simulator/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {

	// init redis and gin
	rdb := repository.NewRedisClient()
	r := gin.Default()

	// init services
	walletService := services.NewWalletService(rdb)

	// init handlers
	walletHandler := handlers.NewWalletHandler(walletService)

	// routes
	// health and chaos

	r.GET("/health", handlers.HealthCheck)
	r.POST("/chaos", handlers.Chaos)

	// wallets
	r.POST("/wallets/:wallet_id/stocks/:stock_name", walletHandler.HandleStockOperation)
	r.GET("/wallets/:wallet_id", walletHandler.GetWallet)
	r.GET("/wallets/:wallet_id/stocks/:stock_name", walletHandler.GetWalletStockQuantity)

	// Bank
	r.GET("/stocks", walletHandler.GetBankState)
	r.POST("/stocks", walletHandler.SetBankState)

	// Log
	r.GET("/log", walletHandler.GetAuditLog)

	r.Run(":8080")

}
