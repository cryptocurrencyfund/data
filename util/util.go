package util

import (
	"encoding/binary"
	"fmt"
	"log"
	"time"

	coinMarketCap "github.com/cryptocurrencyfund/go-coinmarketcap"
)

// FetchPrices Fetches price from API
func FetchPrices(amount int) []coinMarketCap.Coin {
	// Get top 10 coins
	top100, err := coinMarketCap.GetAllCoinDataSorted(amount)
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

// TimeNow time now in unix
func TimeNow() int64 {
	return time.Now().Unix()
}

// ParseTime parse time from unix back to Time object
func ParseTime(u int64) time.Time {
	return time.Unix(u, 0)
}

// Int64ToByteArr int 64 to byte arr
func Int64ToByteArr(i int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	return b
}

// ByteArrToInt64 byte arr to int 64
func ByteArrToInt64(b []byte) int64 {
	i := binary.LittleEndian.Uint64(b)
	return int64(i)
}
