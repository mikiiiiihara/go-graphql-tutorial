package resolver

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/graphql-go/graphql"
)

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