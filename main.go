package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	bolt "github.com/coreos/bbolt"
	"github.com/cryptocurrencyfund/data/util"
	coinMarketCap "github.com/cryptocurrencyfund/go-coinmarketcap"
)

func usage() {
	fmt.Println("Get variable names")
	fmt.Println("./data [number of seconds]")
	er := errors.New("Wrong arguments")
	fmt.Println(er)
}

func main() {
	db := util.OpenDb()
	defer util.CloseDb(db)

	seconds := os.Args[1]
	fmt.Println("Getting price every " + seconds + " seconds")

	top100 := fetchPrices(db)
	saveJSON(util.DateString(), top100)
}

func fetchPrices(db *bolt.DB) []coinMarketCap.Coin {
	// Get top 10 coins
	top100, err := coinMarketCap.GetAllCoinDataSorted(100)
	if err != nil {
		log.Println(err)
	}

	return top100
}

// Top 100 save
func saveJSON(date string, top100 []coinMarketCap.Coin) {

	b, err := json.Marshal(top100)
	if err != nil {
		fmt.Println("error:", err)
	}
	jsonToDisk(date, b)
}

func jsonToDisk(date string, bytes []byte) {

	// open output file
	fo, err := os.Create("./dbFile/" + date + ".json")
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
