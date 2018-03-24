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
	fmt.Println("all: get data, output to json/csv, generate report all in one, push to github")
	fmt.Println("history: fetch all historical prices for top X number of coins")
	fmt.Println("currencyHistory: fetch historical prices for a specific currency")
	fmt.Println("coinInfo: fetch coin information for a specific currency")
	fmt.Println("allCoinInfo: fetch coin information for all currencies")
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
			topPricesToJSON()
			topPricesToCSV()
			dailyReport()
			time.Sleep(time.Duration(24) * time.Hour)
		}
	case "history":
		crawlAllHistory()
		break
	case "currencyHistory":
		crawlCurrencyHistory("ethereum")
		break
	case "allCoinInfo":
		crawlAllCoinInfo()
		break
	case "coinInfo":
		crawlCaseInfo("ethereum")
		break
	case "json":
		for {
			topPricesToJSON()
			time.Sleep(time.Duration(24) * time.Hour)
		}
	case "csv":
		for {
			topPricesToCSV()
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

func topPricesToJSON() {
	top := util.FetchPrices(topCount)
	dateString := util.DateString()
	util.SaveTopPrices(dateString, top)
	util.SyncGit(dateString)
	fmt.Println("Top JSON saved to file")
}

func dailyReport() {
	top := util.FetchPrices(topCount)
	dateString := util.DateString()
	util.GenerateReport(dateString, top)
	util.SyncGit(dateString)
	fmt.Println("Daily report saved to file")
}
func topPricesToCSV() {
	dateString := util.DateString()
	util.GenerateCsv(dateString)
	util.SyncGit(dateString)
	fmt.Println("Top CSV saved to file")
}

func crawlAllHistory() {
	top := util.FetchPrices(500)
	dateString := util.DateString()
	h := reference.CrawlHistoricalData(dateString, top)
	util.SaveHistorialPrices(h)
}

func crawlCurrencyHistory(currency string) {
	dateString := util.DateString()
	reference.CrawlCurrency(dateString, currency)
}

func crawlAllCoinInfo() {
	top := util.FetchPrices(1500)
	infos := reference.CrawlAllCoinInfo(top)
	util.SaveCoinInfo(infos)
}

func crawlCaseInfo(currency string) {
	reference.CrawlCoinInfo(currency)
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
