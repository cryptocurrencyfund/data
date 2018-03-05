package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cryptocurrencyfund/data/reference"
	"github.com/cryptocurrencyfund/data/util"
)

func usage() {
	fmt.Println("\n\n========== Help ==========")
	fmt.Println("json: 24 hour job to fetch data for fund reports")
	fmt.Println("db [seconds]: * second job to fetch data for top 100 coins, default 5")
	fmt.Println("get [coinName | detailed]: get entries from database")
	fmt.Println("help: this screen")
}

func main() {

	switch os.Args[1] {
	case "all":
		for {
			topJSON()
			topCSV()
			dailyReport()
			time.Sleep(time.Duration(24) * time.Hour)
		}
	case "history":
		crawlAllHistory()
	case "currencyHistory":
		crawlCurrencyHistory("ethereum")
	case "json":
		for {
			topJSON()
			time.Sleep(time.Duration(24) * time.Hour)
		}
	case "csv":
		for {
			topCSV()
			time.Sleep(time.Duration(24) * time.Hour)
		}
	case "report":
		for {
			dailyReport()
			time.Sleep(time.Duration(24) * time.Hour)
		}
	case "db":
		seconds := 5
		if len(os.Args) > 2 {
			seconds, _ = strconv.Atoi(os.Args[2])
		}
		fmt.Printf("Getting price every %d seconds\n", seconds)
		topDB(seconds)
		break
	case "get":
		coinType := "Ethereum"
		detailed := false
		if len(os.Args) > 2 {
			coinType = os.Args[2]
		}
		if len(os.Args) > 3 {
			detailed, _ = strconv.ParseBool(os.Args[3])
		}
		fmt.Printf("Getting price for coin type: %s\n", coinType)
		getCoin(coinType, detailed)
		break
	case "help":
		usage()
		break
	default:
		usage()
		break
	}

}

func topJSON() {
	top100 := util.FetchPrices(100)
	dateString := util.DateString()
	util.SaveJSONToFile(dateString, top100)
	util.GenerateReport(dateString, top100)
	util.SyncGit(dateString)
}

func dailyReport() {
	top100 := util.FetchPrices(100)
	dateString := util.DateString()
	util.GenerateReport(dateString, top100)
	util.SyncGit(dateString)
}
func topCSV() {
	dateString := util.DateString()
	util.GenerateCsv(dateString)
	util.SyncGit(dateString)
}

func crawlAllHistory() {
	top := util.FetchPrices(100)
	dateString := util.DateString()
	reference.CrawlHistoricalData(dateString, top)
}

func crawlCurrencyHistory(currency string) {
	dateString := util.DateString()
	reference.CrawlCurrency(dateString, currency)
}

func topDB(seconds int) {
	for {
		db := util.OpenDb()
		currentTs := util.TimeNow()
		top100 := util.FetchPrices(100)
		for i := 0; i < len(top100); i++ {
			util.UpdateCoin(db, currentTs, top100[i])
		}
		util.CloseDb(db)
		time.Sleep(time.Duration(seconds) * time.Second)
	}
}

func getCoin(coinName string, detailed bool) {
	coinName = strings.ToLower(coinName)
	db := util.OpenDb()
	defer util.CloseDb(db)

	coinEntries := util.GetCoin(db, coinName)

	if detailed {
		fmt.Println("Coin Entries: " + coinName)
		fmt.Println(coinEntries)
	}
}
