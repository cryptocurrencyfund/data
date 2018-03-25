package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cryptocurrencyfund/data/reference"
	coinMarketCap "github.com/cryptocurrencyfund/go-coinmarketcap"
)

// SaveTopPrices Top 100 save
func SaveTopPrices(date string, top []coinMarketCap.Coin) {

	b, err := json.Marshal(top)
	if err != nil {
		fmt.Println("error:", err)
	}
	jsonToDisk(date, b)
}

func jsonToDisk(date string, bytes []byte) {

	// open output file
	fo, err := os.Create("data/" + YearString() + "/" + date + ".json")
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

// SaveCoinInfo SaveCoinInfo
func SaveCoinInfo(infos []reference.CoinInfo) {
	bytes, err := json.Marshal(infos)
	if err != nil {
		fmt.Println("error:", err)
	}

	// write json to disk
	fo, err := os.Create("reference/coinInfo/all.json")
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

// SaveHistorialPrices SaveHistorialPrices
func SaveHistorialPrices(p reference.HistorialPrices) {

	bytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println("error:", err)
	}

	// write json to disk
	fo, err := os.Create("reference/historical/all.json")
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

// GetHistoricalPrices GetHistoricalPrices
func GetHistoricalPrices() (p reference.HistorialPrices) {
	file, e := ioutil.ReadFile("reference/historical/all.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	json.Unmarshal(file, &p)
	return
}
