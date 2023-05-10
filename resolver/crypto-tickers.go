package resolver

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/graphql-go/graphql"
)

type Ticker struct {
	ProductCode      string  `json:"product_code"`
	State            string  `json:"state"`
	Timestamp        string  `json:"timestamp"`
	TickID           int     `json:"tick_id"`
	BestBid          float64 `json:"best_bid"`
	BestAsk          float64 `json:"best_ask"`
	BestBidSize      float64 `json:"best_bid_size"`
	BestAskSize      float64 `json:"best_ask_size"`
	TotalBidDepth    float64 `json:"total_bid_depth"`
	TotalAskDepth    float64 `json:"total_ask_depth"`
	MarketBidSize    float64 `json:"market_bid_size"`
	MarketAskSize    float64 `json:"market_ask_size"`
	Ltp              float64 `json:"ltp"`
	Volume           float64 `json:"volume"`
	VolumeByProduct  float64 `json:"volume_by_product"`
}

// 仮想通貨取得
func (r *Resolver) GetCryptoTickers(p graphql.ResolveParams) (interface{}, error) {
	// Get the API URL from the environment variable
	apiURL := os.Getenv("CRYPTO_API_URL")
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ticker Ticker
	err = json.NewDecoder(resp.Body).Decode(&ticker)
	if err != nil {
		return nil, err
	}
	return ticker, nil
}