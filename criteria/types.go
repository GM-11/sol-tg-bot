package criteria

type Trade struct {
	BuyPrice  float64
	SellPrice float64
	Quantity  float64
	Timestamp int64 // Unix timestamp
}

type WalletMetrics struct {
	TotalPnl          float64
	Winrate           float64
	Roi               float64
	AvgPnl            float64
	MaxUnrealizedLoss float64
	MinWinrate        float64
	Min2x             float64
	Min5x             float64
	MinPnlr           float64
	Min30dPnl         float64
	Min7dPnl          float64
	MinAvgPnl         float64
	MinMaxWin         float64
	MaxTop5VsTotal    float64
	HoldTime          int64
	BuySize           float64
}
