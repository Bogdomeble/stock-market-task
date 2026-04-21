
# Stock Market Simulator

Simple REST service simulating a basic stock market with wallets and a central bank.

## Tech stack

* Go (Gin framework)
* Redis (Database & Lua scripts for atomicity)
* Docker & Docker Compose
* Nginx (Load balancing / HA)

## Features

* Buy / Sell stocks (Atomic operations)
* Wallet state management
* Bank as single liquidity provider
* High availability (Multiple Go instances)
* Chaos endpoint (Kills one instance, load balancer handles the failover)
* Audit log tracking

## Requirements

* Docker
* Docker Compose
* Bash (for startup script)

## Run

Solution can be started using a single command where the parameter is the desired port:

```bash
chmod +x run.sh && ./run.sh 8080
or
PORT=8080 docker-compose up --build
```

Application will be available at: `http://localhost:8080`

---

## API

### Wallet Operations

* **Buy / Sell stock**
  `POST /wallets/{wallet_id}/stocks/{stock_name}`
  *Body:* `{"type": "buy" | "sell"}`

* **Get wallet state**
  `GET /wallets/{wallet_id}`

* **Get single stock quantity**
  `GET /wallets/{wallet_id}/stocks/{stock_name}`

### Bank Operations

* **Get bank state**
  `GET /stocks`

* **Set bank state**
  `POST /stocks`
  *Body:* `{"stocks":[{"name":"AAPL", "quantity":99}, {"name":"GOOG", "quantity":1}]}`

### System & Audit

* **Get audit log**
  `GET /log`
  *Returns up to 10,000 successful operations.*

* **Health check**
  `GET /health`

* **Chaos (kill instance)**
  `POST /chaos`

---

## Notes

* Stock price is fixed to 1.
* No wallet balance or funds tracking.
* No order book.
* Concurrency and race conditions are handled via Redis Lua scripts.
* Highly Available (HA): Killing one backend instance via `/chaos` does not stop the system; Nginx routes traffic to the surviving instance.
