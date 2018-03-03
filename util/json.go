package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	coinMarketCap "github.com/cryptocurrencyfund/go-coinmarketcap"
)

// SaveJSONToFile Top 100 save
func SaveJSONToFile(date string, top []coinMarketCap.Coin) {

	b, err := json.Marshal(top)
	if err != nil {
		fmt.Println("error:", err)
	}
	jsonToDisk(date, b)
}

func jsonToDisk(date string, bytes []byte) {

	// open output file
	fo, err := os.Create("./data/" + date + ".json")
	if err != nil {
		panic(err)
	}

	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	// make a write buffer
	w := bufio.NewWriter(fo)
	w.Write(bytes)
	w.Flush()
}
