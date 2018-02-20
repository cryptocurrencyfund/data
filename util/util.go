package util

import (
	"fmt"
	"log"
	"time"

	bolt "github.com/coreos/bbolt"
	coinMarketCap "github.com/cryptocurrencyfund/go-coinmarketcap"
)

// FetchPrices Fetches price from API
func FetchPrices(db *bolt.DB) []coinMarketCap.Coin {
	// Get top 10 coins
	top100, err := coinMarketCap.GetAllCoinDataSorted(100)
	if err != nil {
		log.Println(err)
	}

	return top100
}

// DateString export date as a string
func DateString() string {
	y, m, d := time.Now().Date()
	mStr := fmt.Sprintf("%d", m)
	dStr := fmt.Sprintf("%d", d)
	if m < 10 {
		mStr = fmt.Sprintf("0%d", m)
	}
	if d < 10 {
		dStr = fmt.Sprintf("0%d", d)
	}
	return fmt.Sprintf("%d-%s-%s", y, mStr, dStr)

}

// TimeNow time now
func TimeNow() time.Time {
	return time.Now()
}
