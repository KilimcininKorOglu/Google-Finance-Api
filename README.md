# Google Finance API

Google Finance'ın dahili RPC endpoint'ini sarmalayan, sıfır bağımlılıklı Go REST API.

API anahtarı gerektirmez. Tek bir HTTP isteği ile fiyat, şirket bilgisi, grafik, haber, finansal tablo ve piyasa verileri getirir.

[English README](README.en.md)

## Kurulum

```bash
go build -o google-finance-api ./cmd/server
```

## Çalıştırma

```bash
./google-finance-api
```

Varsayılan port: `8080`. Değiştirmek için:

```bash
PORT=3000 ./google-finance-api
```

## Docker

```bash
docker compose up -d
```

Durdurma ve loglar:

```bash
docker compose logs -f
docker compose down
```

## Ticker Formatı

| Tür     | Format        | Örnek            |
|---------|---------------|------------------|
| Hisse   | SEMBOL:BORSA  | THYAO:IST        |
| Endeks  | .SEMBOL:BORSA | .DJI:INDEXDJX    |
| Kripto  | BAZ-KARŞI     | BTC-USD          |
| Döviz   | BAZ-KARŞI     | EUR-USD          |
| ETF     | SEMBOL:BORSA  | SPY:NYSEARCA     |

## API Endpoint'leri

### Ticker Bazlı

```
GET /v1/quote/{ticker}
GET /v1/company/{ticker}
GET /v1/chart/{ticker}?range=1M
GET /v1/news/{ticker}
GET /v1/financials/{ticker}?type=quarterly
GET /v1/related/{ticker}
GET /v1/full/{ticker}?range=1M
```

### Piyasa

```
GET /v1/market/indices
GET /v1/market/movers?category=most-active&count=10&offset=0
GET /v1/market/trending
GET /v1/market/earnings
GET /v1/market/headlines
```

### Sistem

```
GET /healthz
```

## Örnekler

Hisse fiyatı:

```bash
curl http://localhost:8080/v1/quote/THYAO:IST
```

```json
{
  "ticker": "THYAO",
  "exchange": "IST",
  "name": "Turk Hava Yollari AO",
  "type": "stock",
  "currency": "TRY",
  "timezone": "Europe/Istanbul",
  "price": 325.00,
  "change": 1.50,
  "changePercent": 0.46,
  "previousClose": 323.50
}
```

Şirket bilgisi:

```bash
curl http://localhost:8080/v1/company/GARAN:IST
```

```json
{
  "description": "Garanti BBVA is a Turkish financial services company...",
  "ceo": "Mahmut Akten",
  "employees": 23152,
  "marketCap": 579600000000,
  "open": 138.80,
  "high": 139.90,
  "low": 136.40,
  "fiftyTwoWeekHigh": 169.70,
  "fiftyTwoWeekLow": 98.75,
  "peRatio": 5.28,
  "volume": 28709397,
  "sector": "Bank"
}
```

Tam veri (fiyat + şirket + grafik + haber):

```bash
curl http://localhost:8080/v1/full/ASELS:IST?range=1Y
```

Finansal tablolar (yıllık):

```bash
curl http://localhost:8080/v1/financials/KCHOL:IST?type=annual
```

```json
[
  {
    "fiscalEnd": "2025",
    "isAnnual": true,
    "currency": "TRY",
    "revenue": 2757295000000,
    "netIncome": 22001000000,
    "epsDiluted": 8.68,
    "peRatio": 2.84
  }
]
```

Kripto:

```bash
curl http://localhost:8080/v1/quote/BTC-USD
```

Piyasa endeksleri:

```bash
curl http://localhost:8080/v1/market/indices
```

## Grafik Aralıkları

| Değer | Anlam        |
|-------|--------------|
| 1D    | 1 gün        |
| 5D    | 5 gün        |
| 1M    | 1 ay         |
| 6M    | 6 ay         |
| YTD   | Yıl başı     |
| 1Y    | 1 yıl        |
| 5Y    | 5 yıl        |
| MAX   | Tüm zamanlar |

## Yapı

```
cmd/server/main.go          Giriş noktası, graceful shutdown
internal/gfrpc/client.go    Google Finance RPC istemcisi
internal/gfrpc/codec.go     batchexecute istek/yanıt kodlayıcı
internal/gfrpc/tuple.go     Ticker tuple dönüştürme
internal/gfrpc/methods.go   RPC metot tanımları
internal/decode/            Pozisyonel dizi çözücüleri
internal/api/               HTTP sunucu, handler, middleware
internal/models/            Veri modelleri
```

## Lisans

MIT
