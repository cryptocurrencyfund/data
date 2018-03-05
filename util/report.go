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
	toc(filename)
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
	str := "### Table of contents\n"
	str += "1. Price Changes\n"
	str += "	1. [Price Change Winners (Top 100)](#price-change-winners-top-100)\n"
	str += "	2. [Price Change Losers (Top 100)](#price-change-losers-top-100)\n"
	str += "	1. [Price Change Winners (Top 100)](#price-change-winners-top-300)\n"
	str += "	2. [Price Change Losers (Top 300)](#price-change-losers-top-300)\n"
	str += "2. Volume\n"
	str += "	1. [24H Volume Winners(Top 100)](#24h-volume-winners-top-100)\n"
	str += "	2. [24H Volume Losers(Top 100)](#24h-volume-losers-top-100)\n"
	str += "	1. [24H Volume Winners(Top 300)](#24h-volume-winners-top-300)\n"
	str += "	2. [24H Volume Losers(Top 300)](#24h-volume-losers-top-300)\n"

	if _, err = f.WriteString(str); err != nil {
		panic(err)
	}
}

func priceChangeMd(filename string, top []coinMarketCap.Coin) {
	top100 := make([]coinMarketCap.Coin, 100)
	copy(top100, top[:100])
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	sort.Slice(top, func(i, j int) bool {
		return top[i].PercentChange24h > top[j].PercentChange24h
	})
	sort.Slice(top100, func(i, j int) bool {
		return top100[i].PercentChange24h > top100[j].PercentChange24h
	})
	winners := top[:10]
	losers := top[len(top)-10:]
	winners100 := top100[:10]
	losers100 := top100[len(top100)-10:]

	str := "\n#### Price Change Winners (Top 100)\n"
	for _, w := range winners100 {
		str += w.MarkdownPrice() + "\n"
	}
	str += "\n#### Price Change Losers (Top 100)\n"
	for _, l := range losers100 {
		str += l.MarkdownPrice() + "\n"
	}
	str += "\n#### Price Change Winners (Top 300)\n"
	for _, w := range winners {
		str += w.MarkdownPrice() + "\n"
	}
	str += "\n#### Price Change Losers (Top 300)\n"
	for _, l := range losers {
		str += l.MarkdownPrice() + "\n"
	}

	if _, err = f.WriteString(str); err != nil {
		panic(err)
	}
}

func volMd(filename string, top []coinMarketCap.Coin) {
	top100 := make([]coinMarketCap.Coin, 100)
	copy(top100, top[:100])
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	sort.Slice(top, func(i, j int) bool {
		return top[i].Usd24hVolume > top[j].Usd24hVolume
	})
	sort.Slice(top100, func(i, j int) bool {
		return top100[i].Usd24hVolume > top100[j].Usd24hVolume
	})
	volumeWinners := top[:10]
	volumeLosers := top[len(top)-10:]
	volumeWinners100 := top100[:10]
	volumeLosers100 := top100[len(top100)-10:]
	str := "\n#### 24H Volume Winners (Top 100)\n"
	for _, w := range volumeWinners {
		str += w.MarkdownVolume() + "\n"
	}
	str += "\n#### 24H Volume Losers (Top 100)\n"
	for _, l := range volumeLosers {
		str += l.MarkdownVolume() + "\n"
	}
	str += "\n#### 24H Volume Winners (Top 300)\n"
	for _, w := range volumeWinners100 {
		str += w.MarkdownVolume() + "\n"
	}
	str += "\n#### 24H Volume Losers (Top 300)\n"
	for _, l := range volumeLosers100 {
		str += l.MarkdownVolume() + "\n"
	}

	if _, err = f.WriteString(str); err != nil {
		panic(err)
	}
}

// GeneratePortfolio generate portfolio
func GeneratePortfolio(p *Portfolio) {

}
