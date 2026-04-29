# Google Finance API

Zero-dependency Go REST API that wraps Google Finance's internal RPC endpoint.

No API key required. Fetches price quotes, company info, charts, news, financial statements and market data in a single HTTP request.

[Turkce README](README.md)

## Demo

**https://finance.hermestech.uk**

Hacker-themed terminal UI with live price feed, interactive API explorer, and OpenAPI documentation.

## Installation

```bash
go build -o google-finance-api ./cmd/server
```

## Usage

```bash
./google-finance-api
```

Default port: `8080`. To change:

```bash
PORT=3000 ./google-finance-api
```

## Docker

```bash
docker compose up -d
```

Logs and shutdown:

```bash
docker compose logs -f
docker compose down
```

## Ticker Format

| Type   | Format           | Example          |
|--------|------------------|------------------|
| Stock  | SYMBOL:EXCHANGE  | GOOGL:NASDAQ     |
| Index  | .SYMBOL:EXCHANGE | .DJI:INDEXDJX    |
| Crypto | BASE-QUOTE       | BTC-USD          |
| Forex  | BASE-QUOTE       | EUR-USD          |
| ETF    | SYMBOL:EXCHANGE  | SPY:NYSEARCA     |

## API Endpoints

### Ticker-Based

```
GET /v1/quote/{ticker}
GET /v1/company/{ticker}
GET /v1/chart/{ticker}?range=1M
GET /v1/news/{ticker}
GET /v1/financials/{ticker}?type=quarterly
GET /v1/related/{ticker}
GET /v1/full/{ticker}?range=1M
```

### Market

```
GET /v1/market/indices
GET /v1/market/movers?category=most-active&count=10&offset=0
GET /v1/market/trending
GET /v1/market/earnings
GET /v1/market/headlines
```

### Live Data

```
GET /v1/live              SSE live price stream (15 second interval)
GET /v1/live/snapshot     Instant price JSON
```

Live feed tracks 8 tickers: GOOGL, AAPL, MSFT, BTC-USD, THYAO:IST, USD-TRY, EUR-TRY, EUR-USD.

### System

```
GET /healthz              Health check + version info
GET /openapi.json         OpenAPI 3.1 specification
```

## Examples

Stock quote:

```bash
curl https://finance.hermestech.uk/v1/quote/GOOGL:NASDAQ
```

```json
{
  "ticker": "GOOGL",
  "exchange": "NASDAQ",
  "name": "Alphabet Inc Class A",
  "type": "stock",
  "currency": "USD",
  "timezone": "America/New_York",
  "price": 344.40,
  "change": 5.51,
  "changePercent": 1.63,
  "previousClose": 338.89,
  "afterHours": {
    "price": 343.59,
    "change": -0.81,
    "changePercent": -0.24
  }
}
```

Company info:

```bash
curl https://finance.hermestech.uk/v1/company/AAPL:NASDAQ
```

```json
{
  "description": "Apple Inc. is an American multinational technology company...",
  "ceo": "Tim Cook",
  "employees": 166000,
  "marketCap": 3979467061957,
  "open": 272.76,
  "high": 273.06,
  "low": 269.65,
  "fiftyTwoWeekHigh": 288.61,
  "fiftyTwoWeekLow": 193.25,
  "peRatio": 34.29,
  "volume": 41339643,
  "sector": "Computers, Peripherals, and Software"
}
```

Full data (quote + company + chart + news):

```bash
curl https://finance.hermestech.uk/v1/full/AAPL:NASDAQ?range=1Y
```

Financial statements (annual):

```bash
curl https://finance.hermestech.uk/v1/financials/MSFT:NASDAQ?type=annual
```

```json
[
  {
    "fiscalEnd": "2025",
    "isAnnual": true,
    "currency": "USD",
    "revenue": 281724000000,
    "netIncome": 101832000000,
    "eps": 13.64,
    "peRatio": 19.55
  }
]
```

Crypto:

```bash
curl https://finance.hermestech.uk/v1/quote/BTC-USD
```

```json
{
  "ticker": "BTC-USD",
  "name": "Bitcoin (BTC / USD)",
  "type": "crypto",
  "price": 74039.75,
  "change": -142.28,
  "changePercent": -0.19,
  "previousClose": 74182.03
}
```

Market indices:

```bash
curl https://finance.hermestech.uk/v1/market/indices
```

## Chart Ranges

| Value | Period       |
|-------|--------------|
| 1D    | 1 day        |
| 5D    | 5 days       |
| 1M    | 1 month      |
| 6M    | 6 months     |
| YTD   | Year to date |
| 1Y    | 1 year       |
| 5Y    | 5 years      |
| MAX   | All time     |

## Project Structure

```
cmd/server/main.go              Entry point, graceful shutdown
internal/gfrpc/client.go        Google Finance RPC client
internal/gfrpc/codec.go         batchexecute request/response codec
internal/gfrpc/tuple.go         Ticker tuple conversion
internal/gfrpc/methods.go       RPC method definitions
internal/decode/                Positional array decoders
internal/api/server.go          Route registration (Go 1.22+ method patterns)
internal/api/handlers.go        Ticker-based handlers
internal/api/handlers_market.go Market endpoint handlers
internal/api/handlers_sse.go    SSE live price stream
internal/api/handlers_web.go    Landing page and OpenAPI serving
internal/api/middleware.go      CORS, logging, recovery
internal/models/                Data models
web/index.html                  Terminal-themed landing page
web/openapi.json                OpenAPI 3.1 specification
```

## License

MIT
