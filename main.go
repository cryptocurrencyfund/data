package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cryptocurrencyfund/data/report"
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
	case "json":
		top100JSON()
		break
	case "db":
		seconds := 5
		if len(os.Args) > 2 {
			seconds, _ = strconv.Atoi(os.Args[2])
		}
		fmt.Printf("Getting price every %d seconds\n", seconds)
		top100DB(seconds)
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

func top100JSON() {
	for {
		top100 := util.FetchPrices(100)
		dateString := util.DateString()
		util.SaveJSONToFile(dateString, top100)
		util.SyncGit(dateString)
		report.Generate(dateString, top100)
		time.Sleep(time.Duration(24) * time.Hour)
	}
}

func top100DB(seconds int) {
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
