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
	fmt.Println("\n\n[HELP]")
	fmt.Println("\n\n========== Data ==========")
	fmt.Println("history: fetch all historical prices for top X number of coins")
	fmt.Println("currencyHistory: fetch historical prices for a specific currency")
	fmt.Println("coinInfo: fetch coin information for a specific currency")
	fmt.Println("allCoinInfo: fetch coin information for all currencies")
	fmt.Println("\n========== Charts ==========")
	fmt.Println("currencyCharts: currency charts with price/market cap")
	fmt.Println("comparisonCharts: comparing 2+ currencies charts")
	fmt.Println("themeCharts: theme-based charts")
	fmt.Println("portfolioCharts: portfolio charts")
	fmt.Println("\n========== Cron Jobs ==========")
	fmt.Println("all: get data, output to json/csv, generate report all in one, push to github")
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
	case "currencyCharts":
		currencyCharts()
		break
	case "comparisonCharts":
		comparisonCharts()
		break
	case "themeCharts":
		themeCharts()
		break
	case "portfolioCharts":
		portfolioCharts()
		break
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

func currencyCharts() {
	priceMap := util.GetHistoricalPrices()
	for k := range priceMap {
		currency := k
		util.DrawCurrencyChart(currency)
	}
}

func comparisonCharts() {
	c1 := "bitcoin"
	c2 := "ethereum"
	util.DrawComparisonChart(c1, c2)
}

func themeCharts() {
	util.DrawThemeChart(1000.00, "2018-01-01", "large_cap", "bitcoin", "ethereum", "ripple", "litecoin", "bitcoin-cash", "zcash", "monero")
	util.DrawThemeChart(1000.00, "2018-01-01", "protocol", "ethereum", "omisego", "stellar", "neo", "icon", "zilliqa", "trinity-network-credit", "raiden-network-token")
	util.DrawThemeChart(1000.00, "2018-01-01", "exchanges", "binance-coin", "augur", "gnosis-gno", "0x", "kyber-network", "bancor", "airswap", "republic-protocol")
	util.DrawThemeChart(1000.00, "2018-01-01", "stables", "maker", "trust", "digixdao", "dai", "tether")
	util.DrawThemeChart(1000.00, "2018-01-01", "utility", "basic-attention-token", "golem-network-tokens", "request-network", "funfair", "quantstamp")
	util.DrawThemeChart(1000.00, "2018-01-01", "asia", "neo", "tron", "trinity-network-credit", "icon", "ontology", "storm", "bytom", "vechain", "nebulas-token")
}

func portfolioCharts() {
	portfolio := []util.DumbPortfolio{
		util.DumbPortfolio{
			Currency: "bitcoin",
			Weight:   0.3,
		},
		util.DumbPortfolio{
			Currency: "ethereum",
			Weight:   0.3,
		},
		util.DumbPortfolio{
			Currency: "litecoin",
			Weight:   0.2,
		},
		util.DumbPortfolio{
			Currency: "binance-coin",
			Weight:   0.2,
		},
	}
	util.DrawPortfolioChart(1000.00, "2018-01-01", "casm", portfolio)

	portfolio = []util.DumbPortfolio{
		util.DumbPortfolio{
			Currency: "bitcoin",
			Weight:   0.2,
		},
		util.DumbPortfolio{
			Currency: "ethereum",
			Weight:   0.2,
		},
		util.DumbPortfolio{
			Currency: "litecoin",
			Weight:   0.2,
		},
		util.DumbPortfolio{
			Currency: "zcash",
			Weight:   0.2,
		},
		util.DumbPortfolio{
			Currency: "ripple",
			Weight:   0.2,
		},
	}
	util.DrawPortfolioChart(1000.00, "2018-01-01", "casm2", portfolio)

	portfolio = []util.DumbPortfolio{
		util.DumbPortfolio{
			Currency: "neo",
			Weight:   0.2,
		},
		util.DumbPortfolio{
			Currency: "tron",
			Weight:   0.2,
		},
		util.DumbPortfolio{
			Currency: "vechain",
			Weight:   0.2,
		},
		util.DumbPortfolio{
			Currency: "storm",
			Weight:   0.2,
		},
		util.DumbPortfolio{
			Currency: "nebulas-token",
			Weight:   0.2,
		},
	}
	util.DrawPortfolioChart(1000.00, "2018-01-01", "china", portfolio)

	portfolio = []util.DumbPortfolio{
		util.DumbPortfolio{
			Currency: "binance-coin",
			Weight:   1 / 8,
		},
		util.DumbPortfolio{
			Currency: "augur",
			Weight:   1 / 8,
		},
		util.DumbPortfolio{
			Currency: "gnosis-gno",
			Weight:   1 / 8,
		},
		util.DumbPortfolio{
			Currency: "0x",
			Weight:   1 / 8,
		},
		util.DumbPortfolio{
			Currency: "kyber-network",
			Weight:   1 / 8,
		},
		util.DumbPortfolio{
			Currency: "bancor",
			Weight:   1 / 8,
		},
		util.DumbPortfolio{
			Currency: "airswap",
			Weight:   1 / 8,
		},
		util.DumbPortfolio{
			Currency: "republic-protocol",
			Weight:   1 / 8,
		},
	}
	util.DrawPortfolioChart(1000.00, "2018-01-01", "exchange", portfolio)
}
