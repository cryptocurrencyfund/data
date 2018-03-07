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
func CrawlCoinInfo(currency string) {
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
	writer.Write([]string{"Date", "Open", "High", "Low", "Close", "Volume", "Market Cap"})

	// Instantiate default collector
	c := colly.NewCollector()

	c.OnHTML(".list-unstyled li", func(e *colly.HTMLElement) {
		// fmt.Println("\n\n" + e.ChildText(""))
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
	url := fmt.Sprintf("https://coinmarketcap.com/currencies/%s",
		currency)
	c.Visit(url)
}
