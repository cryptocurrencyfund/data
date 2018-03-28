package util

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/wcharczuk/go-chart"
)

// DrawCurrencyChart DrawCurrencyChart
func DrawCurrencyChart(currency string) {

	timestamps, coinPrices, marketCaps := getCoinPrices(currency)

	priceSeries := chart.TimeSeries{
		Name: "Price ($)",
		Style: chart.Style{
			Show:        true,
			StrokeColor: chart.GetDefaultColor(0),
		},
		XValues: timestamps,
		YValues: coinPrices,
	}

	marketCapSeries := chart.TimeSeries{
		Name: "Market Cap ($)",
		Style: chart.Style{
			Show:        true,
			StrokeColor: chart.GetAlternateColor(0),
		},
		YAxis:   chart.YAxisSecondary,
		XValues: timestamps,
		YValues: marketCaps,
	}

	graph := chart.Chart{
		Width:  4096,
		Height: 500,
		XAxis: chart.XAxis{
			Style:        chart.Style{Show: true},
			TickPosition: chart.TickPositionBetweenTicks,
		},
		YAxis: chart.YAxis{
			Name:      "Price ($)",
			NameStyle: chart.StyleShow(),
			Style:     chart.Style{Show: true},
		},
		YAxisSecondary: chart.YAxis{
			Name:      "MarketCap ($)",
			NameStyle: chart.StyleShow(),
			Style:     chart.Style{Show: true},
		},
		Series: []chart.Series{
			priceSeries,
			marketCapSeries,
		},
	}

	collector := &chart.ImageWriter{}
	graph.Render(chart.PNG, collector)

	image, err := collector.Image()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Final Image: %dx%d\n", image.Bounds().Size().X, image.Bounds().Size().Y)

	outputFile, err := os.Create("charts/currency/" + currency + ".jpg")
	if err != nil {
		// Handle error
	}
	png.Encode(outputFile, image)

	outputFile.Close()
}

// DrawComparisonChart DrawComparisonChart
func DrawComparisonChart(currencies ...string) {

	var allTs [][]time.Time
	var allPrices [][]float64
	for _, v := range currencies {
		ts, p, _ := getCoinPrices(v)
		allTs = append(allTs, ts)
		allPrices = append(allPrices, p)
	}

	var allPriceSeries []chart.Series

	for k, v := range allTs {
		series := chart.TimeSeries{
			Name: currencies[k],
			Style: chart.Style{
				Show:        true,
				StrokeColor: chart.GetDefaultColor(k),
			},
			XValues: v,
			YValues: allPrices[k],
		}
		allPriceSeries = append(allPriceSeries, series)
	}

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
		Series: allPriceSeries,
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	collector := &chart.ImageWriter{}
	graph.Render(chart.PNG, collector)

	image, err := collector.Image()
	if err != nil {
		log.Fatal(err)
	}

	outputFile, err := os.Create("charts/comparison/" + currencies[0] + ".jpg")
	if err != nil {
		// Handle error
	}
	png.Encode(outputFile, image)

	outputFile.Close()
}

// DrawThemeChart DrawThemeChart
func DrawThemeChart(investAmount float64, investDate string, comparisonName string, currencies ...string) {

	var allTs [][]time.Time
	var allPrices [][]float64
	var allValuations [][]float64
	for _, v := range currencies {
		ts, p := getCoinPricesFromDate(v, investDate)
		allTs = append(allTs, ts)
		allPrices = append(allPrices, p)
	}

	for _, v := range allPrices {
		investmentPrice := v[len(v)-1]
		coinsOwned := investAmount / investmentPrice
		var valuations []float64
		for _, p := range v {
			valuations = append(valuations, coinsOwned*p)
		}
		allValuations = append(allValuations, valuations)
	}

	// Print:
	logToFile := "![chart](https://raw.githubusercontent.com/cryptocurrencyfund/data/develop/charts/theme/" + comparisonName + ".jpg)\n\n"
	for c := 0; c < len(currencies); c++ {
		logToFile += "### " + currencies[c] + "\n"
		fmt.Println("[" + currencies[c] + "]")
		for k, v := range allTs[c] {
			current := fmt.Sprintf("Date bought: %s - Current Date: %s - Valuation: %.2f - Price: %.2f \n",
				investDate, v.String(), allValuations[c][k], allPrices[c][k])
			logToFile += "* " + current
			fmt.Print(current)
		}

	}

	var allValuationSeries []chart.Series

	for k, v := range allTs {
		series := chart.TimeSeries{
			Name: currencies[k],
			Style: chart.Style{
				Show:        true,
				StrokeColor: chart.GetDefaultColor(k),
			},
			XValues: v,
			YValues: allValuations[k],
		}
		allValuationSeries = append(allValuationSeries, series)
	}

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

	outputImage, err := os.Create("charts/theme/" + comparisonName + ".jpg")
	outputFile, _ := os.Create("charts/theme/" + comparisonName + ".md")
	if err != nil {
		// Handle error
	}
	png.Encode(outputImage, image)
	outputFile.WriteString(logToFile)
	outputImage.Close()
	outputFile.Close()
}

func getCoinPrices(currencyType string) (timestamps []time.Time, coinPrices []float64, marketCaps []float64) {
	h := GetHistoricalPrices()
	currency := h[currencyType]
	var dates []string
	for _, s := range currency {
		dTime, _ := dateparse.ParseAny(s.Date)
		d := dTime.Format(chart.DefaultDateFormat)
		dates = append(dates, d)
		coinPrices = append(coinPrices, s.Close)
		mkcp, _ := strconv.ParseFloat(s.MarketCap, 64)
		marketCaps = append(marketCaps, mkcp)
	}
	fmt.Printf("%v\n", marketCaps)
	for _, ts := range dates {
		parsed, _ := time.Parse(chart.DefaultDateFormat, ts)
		timestamps = append(timestamps, parsed)
	}

	return
}

func getCoinPricesFromDate(currencyType string, fromDate string) (timestamps []time.Time, coinPrices []float64) {
	h := GetHistoricalPrices()
	currency := h[currencyType]
	var dates []string
	for _, s := range currency {
		dTime, _ := dateparse.ParseAny(s.Date)
		d := dTime.Format(chart.DefaultDateFormat)
		if strings.Compare(d, fromDate) >= 0 {
			dates = append(dates, d)
			coinPrices = append(coinPrices, s.Close)
		}
	}
	for _, ts := range dates {
		parsed, _ := time.Parse(chart.DefaultDateFormat, ts)
		timestamps = append(timestamps, parsed)
	}

	return
}

// smaSeries := chart.SMASeries{
// 	Style: chart.Style{
// 		Show:            true,
// 		StrokeColor:     drawing.ColorRed,
// 		StrokeDashArray: []float64{5.0, 5.0},
// 	},
// 	InnerSeries: priceSeries,
// }

// bbSeries := &chart.BollingerBandsSeries{
// 	Style: chart.Style{
// 		Show:        true,
// 		StrokeColor: drawing.ColorFromHex("efefef"),
// 		FillColor:   drawing.ColorFromHex("efefef").WithAlpha(64),
// 	},
// 	InnerSeries: priceSeries,
// }
