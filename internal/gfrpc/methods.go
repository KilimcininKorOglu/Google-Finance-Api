package gfrpc

const (
	MethodQuote          = "xh8wxf"
	MethodCompany        = "HqGpWd"
	MethodClassification = "uwlMvd"
	MethodFinancials     = "Pr8h2e"
	MethodChart          = "AiCwsd"
	MethodNews           = "nBEQBc"
	MethodAnalyst        = "o6pODe"
	MethodRelated        = "SICF5d"
	MethodStockContext   = "mKsvE"
	MethodMarketIndices  = "Xhdx2e"
	MethodMarketMovers   = "YtbmEe"
	MethodTrending       = "lvVhof"
	MethodEarnings       = "JFUMjd"
	MethodCategoryStocks = "XqaYg"
	MethodTopHeadline    = "QKZUzd"
)

func QuoteRequest(tuple []any) RPCRequest {
	return RPCRequest{ID: MethodQuote, Params: []any{[]any{tuple}, 1}}
}

func CompanyRequest(tuple []any) RPCRequest {
	return RPCRequest{ID: MethodCompany, Params: []any{[]any{tuple}}}
}

func FinancialsRequest(tuple []any) RPCRequest {
	return RPCRequest{ID: MethodFinancials, Params: []any{[]any{tuple}}}
}

func ChartRequest(tuple []any, mode int) RPCRequest {
	return RPCRequest{ID: MethodChart, Params: []any{[]any{tuple}, mode}}
}

func NewsRequest(isCrypto bool, tuple []any) RPCRequest {
	newsType := 5
	if isCrypto {
		newsType = 6
	}
	return RPCRequest{ID: MethodNews, Params: []any{newsType, 3, []any{tuple}}}
}

func RelatedRequest(tuple []any) RPCRequest {
	return RPCRequest{ID: MethodRelated, Params: []any{tuple, 18}}
}

func MarketIndicesRequest() RPCRequest {
	return RPCRequest{ID: MethodMarketIndices, Params: []any{nil, 1}}
}

func MarketMoversRequest(categories []string, count, offset int) RPCRequest {
	return RPCRequest{ID: MethodMarketMovers, Params: []any{categories, count, offset}}
}

func TrendingRequest() RPCRequest {
	return RPCRequest{ID: MethodTrending, Params: []any{18}}
}

func EarningsRequest() RPCRequest {
	return RPCRequest{ID: MethodEarnings, Params: []any{}}
}

func TopHeadlineRequest() RPCRequest {
	return RPCRequest{ID: MethodTopHeadline, Params: []any{1}}
}

func CategoryStocksRequest(category string, offset int) RPCRequest {
	return RPCRequest{ID: MethodCategoryStocks, Params: []any{category, offset}}
}
