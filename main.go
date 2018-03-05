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

const topCount = 300

func usage() {
	fmt.Println("\n\n========== Help ==========")
	fmt.Println("all: get data, output to json/csv, generate report all in one")
	fmt.Println("history: fetch all historical prices for top X number of coins")
	fmt.Println("currencyHistory: fetch historical prices for a specific currency")
	fmt.Println("json: 24 hour job to fetch data and save as json")
	fmt.Println("csv: 24 hour job to fetch data and save as csv")
	fmt.Println("report: 24 hour job to generate a daily report")
	fmt.Println("db [seconds]: * second job to fetch data for top X number of coins, default 5")
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
	top := util.FetchPrices(topCount)
	dateString := util.DateString()
	util.SaveJSONToFile(dateString, top)
	util.GenerateReport(dateString, top)
	util.SyncGit(dateString)
}

func dailyReport() {
	top := util.FetchPrices(topCount)
	dateString := util.DateString()
	util.GenerateReport(dateString, top)
	util.SyncGit(dateString)
}
func topCSV() {
	dateString := util.DateString()
	util.GenerateCsv(dateString)
	util.SyncGit(dateString)
}

func crawlAllHistory() {
	top := util.FetchPrices(500)
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
		top := util.FetchPrices(topCount)
		for i := 0; i < len(top); i++ {
			util.UpdateCoin(db, currentTs, top[i])
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
