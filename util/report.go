package util

import (
	"bufio"
	"os"
	"sort"

	coinMarketCap "github.com/cryptocurrencyfund/go-coinmarketcap"
)

// GenerateReport Generate daily report
func GenerateReport(dateString string, top []coinMarketCap.Coin) {
	filename := "report/" + YearString() + "/" + dateString + ".md"
	createMarkDown(dateString, filename)
	priceChangeMd(filename, top)
	volMd(filename, top)
}

func createMarkDown(date string, filename string) {

	// open output file
	fo, err := os.Create(filename)
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
	w.WriteString("## " + date + "\n")
	w.Flush()
}

func toc(filename string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	str := "### Table of contents"
	str += "1. Price Changes"
	str += "	1. Price Change Winners(#price-change-winners)"
	str += "	2. Price Change Losers(#price-change-losers)"
	str += "2. Volume"
	str += "	1. 24H Volume Winners(#24h-volume-winners)"
	str += "	2. 24H Volume Losers(#24h-volume-losers)"

	if _, err = f.WriteString(str); err != nil {
		panic(err)
	}
}

func priceChangeMd(filename string, top []coinMarketCap.Coin) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	sort.Slice(top, func(i, j int) bool {
		return top[i].PercentChange24h > top[j].PercentChange24h
	})
	winners := top[:10]
	losers := top[len(top)-10:]

	str := "\n#### Price Change Winners\n"
	for _, w := range winners {
		str += w.MarkdownPrice() + "\n"
	}
	str += "\n#### Price Change Losers\n"
	for _, l := range losers {
		str += l.MarkdownPrice() + "\n"
	}

	if _, err = f.WriteString(str); err != nil {
		panic(err)
	}
}

func volMd(filename string, top []coinMarketCap.Coin) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	sort.Slice(top, func(i, j int) bool {
		return top[i].Usd24hVolume > top[j].Usd24hVolume
	})
	volumeWinners := top[:10]
	volumeLosers := top[len(top)-10:]
	str := "\n#### 24H Volume Winners\n"
	for _, w := range volumeWinners {
		str += w.MarkdownVolume() + "\n"
	}
	str += "\n#### 24H Volume Losers\n"
	for _, l := range volumeLosers {
		str += l.MarkdownVolume() + "\n"
	}

	if _, err = f.WriteString(str); err != nil {
		panic(err)
	}
}

// GeneratePortfolio generate portfolio
func GeneratePortfolio(p *Portfolio) {

}
