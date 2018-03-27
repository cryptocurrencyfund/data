package util

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/wcharczuk/go-chart"
)

// DumbPortfolio DumbPortfolio
type DumbPortfolio struct {
	Currency string
	Weight   float64
}

// DrawPortfolioChart DrawPortfolioChart
func DrawPortfolioChart(investAmount float64, investDate string, portfolioName string, portfolio []DumbPortfolio) {

	var allTs [][]time.Time       // [0: {"2018-01-01","2018-01-02"}]
	var allPrices [][]float64     // [0: [1, 2, 3]]
	var allValuations [][]float64 // [priceTotal(0), priceTotal()]

	for _, v := range portfolio {
		ts, p := getCoinPricesFromDate(v.Currency, investDate)
		allTs = append(allTs, ts)
		allPrices = append(allPrices, p)
	}

	for k, v := range allPrices {
		investmentPrice := v[len(v)-1]
		coinsOwned := investAmount * portfolio[k].Weight / investmentPrice
		var valuations []float64
		for _, p := range v {
			valuations = append(valuations, coinsOwned*p)
		}
		allValuations = append(allValuations, valuations) // [0: [1, 3, 5] 1: [2, 4, 6]]
	}
	var compoundValuations []float64 // [0: [1, 3, 5] 1: [2, 4, 6]] => [3, 7, 11]

	for i := 0; i < len(allValuations[0]); i++ {
		sum := 0.0
		for _, v := range allValuations {
			sum += v[i]
		}
		compoundValuations = append(compoundValuations, sum)
	}

	// Print:
	logToFile := "![chart](https://raw.githubusercontent.com/cryptocurrencyfund/data/develop/charts/portfolio/" + portfolioName + ".jpg)\n\n"
	for c := 0; c < len(portfolio); c++ {
		logToFile += "### " + portfolio[c].Currency + "\n"
		fmt.Println("[" + portfolio[c].Currency + "]")
		for k, v := range allTs[c] {
			current := fmt.Sprintf("Date bought: %s - Current Date: %s - Valuation: %.2f - Price: %.2f \n",
				investDate, v.String(), allValuations[c][k], allPrices[c][k])
			logToFile += "* " + current
			fmt.Print(current)
		}

	}
	fmt.Printf("\nCompound: \n%v\n", compoundValuations)
	logToFile += fmt.Sprintf("\nCompound: \n%v\n", compoundValuations)

	var allValuationSeries []chart.Series

	baseSeries := chart.TimeSeries{
		Name: portfolio[0].Currency,
		Style: chart.Style{
			Show:        true,
			StrokeColor: chart.GetDefaultColor(0),
		},
		XValues: allTs[0],
		YValues: allValuations[0],
	}
	portfolioSeries := chart.TimeSeries{
		Name: portfolio[0].Currency,
		Style: chart.Style{
			Show:        true,
			StrokeColor: chart.GetDefaultColor(1),
		},
		XValues: allTs[0],
		YValues: compoundValuations,
	}
	allValuationSeries = append(allValuationSeries, baseSeries)
	allValuationSeries = append(allValuationSeries, portfolioSeries)

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name:         "Time",
			Style:        chart.Style{Show: true},
			TickPosition: chart.TickPositionBetweenTicks,
		},
		YAxis: chart.YAxis{
			Name:  "Price ($)",
			Style: chart.Style{Show: true},
		},
		Background: chart.Style{
			Padding: chart.Box{
				Top:  20,
				Left: 40,
			},
		},
		Series: allValuationSeries,
	}

	graph.Elements = []chart.Renderable{
		chart.LegendLeft(&graph),
	}

	collector := &chart.ImageWriter{}
	graph.Render(chart.PNG, collector)

	image, err := collector.Image()
	if err != nil {
		log.Fatal(err)
	}

	outputImage, err := os.Create("charts/portfolio/" + portfolioName + ".jpg")
	outputFile, _ := os.Create("charts/portfolio/" + portfolioName + ".md")
	if err != nil {
		// Handle error
	}
	png.Encode(outputImage, image)
	outputFile.WriteString(logToFile)
	outputImage.Close()
	outputFile.Close()
}

func getPortfolioPriceFromDate(portfolio DumbPortfolio, fromDate string) {

}
