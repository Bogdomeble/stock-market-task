# Stock Market Simulator

Simple REST service simulating a basic stock market with wallets and a central bank.

## Tech stack

* Go
* Gin
* Redis
* Docker & Docker Compose
* Nginx (load balancing)

## Features

* Buy / Sell stocks
* Wallet state management
* Bank as single liquidity provider
* Basic high availability (multiple instances)
* Chaos endpoint (kills one instance)

## Requirements

* Docker
* Docker Compose

## Run

```bash
PORT=8080 docker compose up --build
```

Application will be available at:

```
http://localhost:8080
```

---

## API

### Buy / Sell stock

```
POST /wallets/{wallet_id}/stocks/{stock_name}
```

Body:

```json
{
  "type": "buy" | "sell"
}
```

---

### Get wallet

```
GET /wallets/{wallet_id}
```

---

### Health check

```
GET /health
```

---

### Chaos (kill instance)

```
POST /chaos
```

---

## Notes

* Stock price is fixed to 1
* No wallet balance tracking
* No order book
* Redis is used as storage
* Multiple instances are load-balanced via Nginx

---

## TODO

* Atomic operations (Lua scripts)
* Audit log endpoint
* Full GET endpoints implementation
* Validation improvements
