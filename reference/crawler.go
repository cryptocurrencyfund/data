package reference

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

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

// CrawlHistoricalData CrawlHistoricalData
func CrawlHistoricalData(date string, top []coinMarketCap.Coin) {
	for _, v := range top {
		currency := strings.ToLower(v.ID)
		CrawlCurrency(date, currency)
	}
}

// CrawlCurrency CrawlCurrency
func CrawlCurrency(date string, currency string) {
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

	c.OnHTML("#historical-data tbody tr", func(e *colly.HTMLElement) {
		writer.Write([]string{
			e.ChildText("td:nth-child(1)"),
			e.ChildText("td:nth-child(2)"),
			e.ChildText("td:nth-child(3)"),
			e.ChildText("td:nth-child(4)"),
			e.ChildText("td:nth-child(5)"),
			e.ChildText("td:nth-child(6)"),
			e.ChildText("td:nth-child(7)"),
		})
	})
	end := strings.Trim(date, "-")
	url := fmt.Sprintf("https://coinmarketcap.com/currencies/%s/historical-data/?start=20110101&end=%s",
		currency,
		end)
	c.Visit(url)
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
