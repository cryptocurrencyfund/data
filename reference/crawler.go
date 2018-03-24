package reference

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"

	coinMarketCap "github.com/cryptocurrencyfund/go-coinmarketcap"
	"github.com/gocolly/colly"
)

// CoinInfo CoinInfo
type CoinInfo struct {
	ID           string `json:"id"`
	Website      string `json:"website"`
	Announcement string `json:"announcement"`
	Explorer     string `json:"explorer"`
	MessageBoard string `json:"messageBoard"`
	Chat         string `json:"chat"`
	Github       string `json:"github"`
}

// HistorialPrice HistorialPrice
type HistorialPrice struct {
	Date      time.Time `json:"date"`
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	Volume    float64   `json:"volume"`
	MarketCap float64   `json:"marketCap"`
}

// HistorialPrices HistoricalPrices
type HistorialPrices map[string][]HistorialPrice

// CrawlHistoricalData CrawlHistoricalData
func CrawlHistoricalData(date string, top []coinMarketCap.Coin) (h HistorialPrices) {
	h = make(HistorialPrices)
	for _, v := range top {
		currency := strings.ToLower(v.ID)
		h[v.ID] = CrawlCurrency(date, currency)
	}

	fmt.Printf("Crawl historical data: %d\n", len(h))
	return
}

// CrawlCurrency CrawlCurrency
func CrawlCurrency(date string, currency string) (arr []HistorialPrice) {
	filename := "reference/historical/" + currency + ".csv"
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", filename, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	writer.Write([]string{"Date", "Open", "High", "Low", "Close", "Volume", "Market Cap"})

	// Instantiate default collector
	c := colly.NewCollector()

	// init historical price array
	var h HistorialPrice

	c.OnHTML("#historical-data tbody tr", func(e *colly.HTMLElement) {
		dateStr := e.ChildText("td:nth-child(1)")
		openStr := e.ChildText("td:nth-child(2)")
		highStr := e.ChildText("td:nth-child(3)")
		lowStr := e.ChildText("td:nth-child(4)")
		closeStr := e.ChildText("td:nth-child(5)")
		volumeStr := e.ChildText("td:nth-child(6)")
		marketCapStr := e.ChildText("td:nth-child(7)")
		h.Date, _ = dateparse.ParseAny(dateStr)
		h.Open, _ = strconv.ParseFloat(openStr, 64)
		h.High, _ = strconv.ParseFloat(highStr, 64)
		h.Low, _ = strconv.ParseFloat(lowStr, 64)
		h.Close, _ = strconv.ParseFloat(closeStr, 64)
		h.Volume, _ = strconv.ParseFloat(volumeStr, 64)
		h.MarketCap, _ = strconv.ParseFloat(marketCapStr, 64)
		arr = append(arr, h)

		writer.Write([]string{
			dateStr,
			openStr,
			highStr,
			lowStr,
			closeStr,
			volumeStr,
			marketCapStr,
		})
	})
	end := strings.Trim(date, "-")
	url := fmt.Sprintf("https://coinmarketcap.com/currencies/%s/historical-data/?start=20110101&end=%s",
		currency,
		end)
	c.Visit(url)
	return
}

// CrawlCoinInfo CrawlCoinInfo
func CrawlCoinInfo(currency string) (info CoinInfo) {
	info.ID = currency
	filename := "reference/coinInfo/" + currency + ".csv"
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", filename, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	writer.Write([]string{"Website", "Announcement", "Explorer", "Message Board", "Chat", "Github"})

	// Instantiate default collector
	c := colly.NewCollector()

	c.OnHTML(".list-unstyled", func(e *colly.HTMLElement) {
		callback := func(index int, innerE *colly.HTMLElement) {
			attrName := innerE.ChildText("a")
			url := innerE.ChildAttr("a", "href")
			fmt.Println("attr " + attrName)
			fmt.Println("url " + url)
			switch attrName {
			case "Website":
				info.Website = url
				break
			case "Announcement":
				info.Announcement = url
				break
			case "Explorer":
				info.Explorer = url
				break
			case "Message Board":
				info.MessageBoard = url
				break
			case "Chat":
				info.Chat = url
				break
			case "Source Code":
				info.Github = url
				break
			}
		}
		e.ForEach("li", callback)
		writer.Write([]string{
			info.Website,
			info.Announcement,
			info.Explorer,
			info.MessageBoard,
			info.Chat,
			info.Github,
		})
	})

	// hit url
	url := fmt.Sprintf("https://coinmarketcap.com/currencies/%s",
		currency)
	c.Visit(url)
	return
}

// CrawlAllCoinInfo CrawlAllCoinInfo
func CrawlAllCoinInfo(top []coinMarketCap.Coin) (infos []CoinInfo) {
	for _, v := range top {
		currency := strings.ToLower(v.ID)
		infos = append(infos, CrawlCoinInfo(currency))
	}
	return
}
