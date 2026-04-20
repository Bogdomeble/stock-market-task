package services

import (
	"fmt"
	"stock-simulator/internal/repository"

	"github.com/redis/go-redis/v9"
)

type WalletService struct {
	rdb *redis.Client
}

func NewWalletService(rdb *redis.Client) *WalletService {
	return &WalletService{rdb: rdb}
}

func (s *WalletService) Buy(walletID, stock string) error {
	bankKey := fmt.Sprintf("bank:%s", stock)
	walletKey := fmt.Sprintf("wallet:%s:%s", walletID, stock)

	bankStock, err := s.rdb.Get(repository.Ctx, bankKey).Int()
	if err != nil || bankStock <= 0 {
		return fmt.Errorf("no stock available")
	}

	s.rdb.Decr(repository.Ctx, bankKey)
	s.rdb.Incr(repository.Ctx, walletKey)

	return nil
}

func (s *WalletService) Sell(walletID, stock string) error {
	bankKey := fmt.Sprintf("bank:%s", stock)
	walletKey := fmt.Sprintf("wallet:%s:%s", walletID, stock)

	walletStock, err := s.rdb.Get(repository.Ctx, walletKey).Int()
	if err != nil || walletStock <= 0 {
		return fmt.Errorf("no stock in wallet")
	}

	s.rdb.Decr(repository.Ctx, walletKey)
	s.rdb.Incr(repository.Ctx, bankKey)

	return nil
}