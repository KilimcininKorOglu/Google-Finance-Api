package models

type Quote struct {
	Ticker        string      `json:"ticker"`
	Exchange      string      `json:"exchange,omitempty"`
	Name          string      `json:"name"`
	Type          string      `json:"type"`
	Currency      string      `json:"currency"`
	Timezone      string      `json:"timezone,omitempty"`
	Price         float64     `json:"price"`
	Change        float64     `json:"change"`
	ChangePercent float64     `json:"changePercent"`
	PreviousClose float64     `json:"previousClose"`
	AfterHours    *AfterHours `json:"afterHours,omitempty"`
}

type AfterHours struct {
	Price         float64 `json:"price"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"changePercent,omitempty"`
}

type CompanyInfo struct {
	Description      string  `json:"description,omitempty"`
	CEO              string  `json:"ceo,omitempty"`
	Employees        int64   `json:"employees,omitempty"`
	MarketCap        float64 `json:"marketCap,omitempty"`
	Open             float64 `json:"open,omitempty"`
	High             float64 `json:"high,omitempty"`
	Low              float64 `json:"low,omitempty"`
	FiftyTwoWeekHigh float64 `json:"fiftyTwoWeekHigh,omitempty"`
	FiftyTwoWeekLow  float64 `json:"fiftyTwoWeekLow,omitempty"`
	PERatio          float64 `json:"peRatio,omitempty"`
	Volume           int64   `json:"volume,omitempty"`
	Sector           string  `json:"sector,omitempty"`
}

type ChartData struct {
	PreviousClose float64      `json:"previousClose"`
	Points        []ChartPoint `json:"points"`
}

type ChartPoint struct {
	Date   string  `json:"date"`
	Price  float64 `json:"price"`
	Volume *int64  `json:"volume,omitempty"`
}

type NewsItem struct {
	Title     string `json:"title"`
	Source    string `json:"source"`
	URL       string `json:"url"`
	Timestamp int64  `json:"timestamp,omitempty"`
}

type FinancialPeriod struct {
	FiscalEnd         string  `json:"fiscalEnd"`
	IsAnnual          bool    `json:"isAnnual"`
	Currency          string  `json:"currency,omitempty"`
	Revenue           float64 `json:"revenue,omitempty"`
	NetIncome         float64 `json:"netIncome,omitempty"`
	EPS               float64 `json:"eps,omitempty"`
	EPSDiluted        float64 `json:"epsDiluted,omitempty"`
	OperatingMargin   float64 `json:"operatingMargin,omitempty"`
	OperatingIncome   float64 `json:"operatingIncome,omitempty"`
	EBITDA            float64 `json:"ebitda,omitempty"`
	SharesOutstanding float64 `json:"sharesOutstanding,omitempty"`
	RevenueGrowthYoY  float64 `json:"revenueGrowthYoY,omitempty"`
	PERatio           float64 `json:"peRatio,omitempty"`
	TotalAssets       float64 `json:"totalAssets,omitempty"`
	TotalLiabilities  float64 `json:"totalLiabilities,omitempty"`
	TotalEquity       float64 `json:"totalEquity,omitempty"`
	OperatingCashFlow float64 `json:"operatingCashFlow,omitempty"`
	ProfitMargin      float64 `json:"profitMargin,omitempty"`
	FreeCashFlow      float64 `json:"freeCashFlow,omitempty"`
	CapEx             float64 `json:"capEx,omitempty"`
}

type MarketIndex struct {
	Ticker        string  `json:"ticker"`
	Name          string  `json:"name"`
	Price         float64 `json:"price"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"changePercent"`
}

type MarketMover struct {
	Ticker        string  `json:"ticker"`
	Name          string  `json:"name"`
	Price         float64 `json:"price"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"changePercent"`
}

type EarningsEvent struct {
	Ticker   string `json:"ticker"`
	Name     string `json:"name"`
	Date     string `json:"date,omitempty"`
	Exchange string `json:"exchange,omitempty"`
}

type Headline struct {
	Title  string `json:"title"`
	URL    string `json:"url"`
	Source string `json:"source,omitempty"`
}

type RelatedStock struct {
	Ticker        string  `json:"ticker"`
	Name          string  `json:"name"`
	Price         float64 `json:"price"`
	Change        float64 `json:"change"`
	ChangePercent float64 `json:"changePercent"`
}

type FullQuote struct {
	Quote   *Quote       `json:"quote"`
	Company *CompanyInfo `json:"company,omitempty"`
	Chart   *ChartData   `json:"chart,omitempty"`
	News    []NewsItem   `json:"news,omitempty"`
}
