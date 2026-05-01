package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// variables for error handling
var (
	ErrStockNotFound     = errors.New("NOT_FOUND")
	ErrInsufficientFunds = errors.New("INSUFFICIENT_FUNDS")
	ErrNoStockInWallet   = errors.New("NO_STOCK_IN_WALLET")
	ErrWalletNotFound    = errors.New("WALLET_NOT_FOUND")
)

type WalletService struct {
	rdb *redis.Client
}

func NewWalletService(rdb *redis.Client) *WalletService { // constructor
	return &WalletService{rdb: rdb}
}

const (
	bankKey = "bank_stocks"
	logKey  = "audit_log"
)

var buyScript = redis.NewScript(`
	local bank_qty = redis.call("HGET", KEYS[1], ARGV[1])
	if not bank_qty then return -1 end
	if tonumber(bank_qty) <= 0 then return -2 end

	redis.call("HINCRBY", KEYS[1], ARGV[1], -1)
	redis.call("HINCRBY", KEYS[2], ARGV[1], 1)
	return 1
`)

var sellScript = redis.NewScript(`
	local bank_exists = redis.call("HEXISTS", KEYS[2], ARGV[1])
	if bank_exists == 0 then return -1 end

	local wallet_qty = redis.call("HGET", KEYS[1], ARGV[1])
	if not wallet_qty or tonumber(wallet_qty) <= 0 then return -2 end

	redis.call("HINCRBY", KEYS[1], ARGV[1], -1)
	redis.call("HINCRBY", KEYS[2], ARGV[1], 1)
	return 1
`)

func (s *WalletService) Buy(walletID, stockName string) error {
	ctx := context.Background()
	walletKey := fmt.Sprintf("wallet:%s", walletID)

	res, err := buyScript.Run(ctx, s.rdb,[]string{bankKey, walletKey}, stockName).Int()
	if err != nil {
		return err
	}

	if res == -1 {
		return ErrStockNotFound
	}
	if res == -2 {
		return ErrInsufficientFunds
	}

	s.logAction("buy", walletID, stockName)
	return nil
}

func (s *WalletService) Sell(walletID, stockName string) error {
	ctx := context.Background()
	walletKey := fmt.Sprintf("wallet:%s", walletID)

	res, err := sellScript.Run(ctx, s.rdb,[]string{walletKey, bankKey}, stockName).Int()
	if err != nil {
		return err
	}

	if res == -1 {
		return ErrStockNotFound
	}
	if res == -2 {
		return ErrNoStockInWallet
	}

	s.logAction("sell", walletID, stockName)
	return nil
}

func (s *WalletService) logAction(opType, walletID, stockName string) {
	entry, _ := json.Marshal(map[string]string{
		"type":       opType,
		"wallet_id":  walletID,
		"stock_name": stockName,
	})

	ctx := context.Background()
	s.rdb.RPush(ctx, logKey, entry)
	// no more than 10k log entries
	s.rdb.LTrim(ctx, logKey, -10000, -1)
}

func (s *WalletService) GetWallet(walletID string) (map[string]int, error) {
	key := "wallet:" + walletID

	// check if the wallet exists for the given walletID
	exists, err := s.rdb.Exists(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}
	if exists == 0 {
		return nil, ErrWalletNotFound
	}

	res, err := s.rdb.HGetAll(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	stocks := make(map[string]int)
	for k, v := range res {
		var qty int
		fmt.Sscanf(v, "%d", &qty)
		if qty > 0 { // do not include stock that have 0 counts on a wallet
			stocks[k] = qty
		}
	}
	return stocks, nil
}

func (s *WalletService) GetWalletStock(walletID, stockName string) (int, error) {
	res, err := s.rdb.HGet(context.Background(), "wallet:"+walletID, stockName).Int()

	if err == redis.Nil {
		return 0, nil
	}

	return res, err
}

func (s *WalletService) GetBankState() (map[string]int, error) {
	res, err := s.rdb.HGetAll(context.Background(), bankKey).Result()
	if err != nil {
		return nil, err
	}

	stocks := make(map[string]int)
	for k, v := range res {
		var qty int
		fmt.Sscanf(v, "%d", &qty)
		stocks[k] = qty
	}
	return stocks, nil
}

func (s *WalletService) SetBankState(stocks map[string]int) error {
	ctx := context.Background()
	pipe := s.rdb.TxPipeline()
	pipe.Del(ctx, bankKey) // clear old state

	if len(stocks) > 0 {
		var args[]interface{}
		for k, v := range stocks {
			args = append(args, k, v)
		}
		pipe.HSet(ctx, bankKey, args...)
	}

	_, err := pipe.Exec(ctx)
	return err
}

func (s *WalletService) GetLog() ([]map[string]string, error) {
	res, err := s.rdb.LRange(context.Background(), logKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	var logs []map[string]string
	for _, item := range res {
		var entry map[string]string
		json.Unmarshal([]byte(item), &entry)
		logs = append(logs, entry)
	}
	if logs == nil {
		logs =[]map[string]string{}
	}
	return logs, nil
}
