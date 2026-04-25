package main

import (
	"stock-simulator/internal/handlers"
	"stock-simulator/internal/repository"
	"stock-simulator/internal/services"

	"github.com/gin-gonic/gin"
	"fmt"
	"net/http"
	"time"
	"os"
	"context"
	"os/signal"
	"syscall"

)

func main() {

	// init redis and gin

	rdb := repository.NewRedisClient()

	handlers.RedisClient = rdb

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

	srv := &http.Server{ //use built-in server instead of gin.run()
	Addr:           ":8080",
	Handler:        r,
	ReadTimeout:    5 * time.Second,
	WriteTimeout:   5 * time.Second,
	IdleTimeout:    30 * time.Second,
	}
		 //sigtem handling, maybe needed maybe not?
		 
		go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("Server Shutdown Error:", err)
	}
}
