package main

import (
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

	fetchPrices(db)
}

func fetchPrices(db *bolt.DB) {
	// Get top 10 coins
	top10, err := coinMarketCap.GetAllCoinDataSorted(10)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(top10)
	}
}
