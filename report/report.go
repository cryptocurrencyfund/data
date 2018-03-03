package report

import (
	"fmt"
	"sort"

	coinMarketCap "github.com/cryptocurrencyfund/go-coinmarketcap"
)

// Report Report
type Report struct {
	Winners       []coinMarketCap.Coin
	Losers        []coinMarketCap.Coin
	VolumeWinners []coinMarketCap.Coin
	VolumeLosers  []coinMarketCap.Coin
}

// Generate generate daily report
func Generate(dateString string, top []coinMarketCap.Coin) {
	var r Report
	sort.Slice(top, func(i, j int) bool {
		return top[i].PercentChange24h > top[j].PercentChange24h
	})
	r.Winners = top[:10]
	r.Losers = top[len(top)-10:]

	sort.Slice(top, func(i, j int) bool {
		return top[i].Usd24hVolume > top[j].Usd24hVolume
	})
	r.VolumeWinners = top[:10]
	r.VolumeLosers = top[len(top)-10:]

	fmt.Printf("Winners: \n%v\n", r.Winners)
	fmt.Printf("Losers: \n%v\n", r.Winners)
	fmt.Printf("Volumn winners: \n%v\n", r.VolumeWinners)
	fmt.Printf("Volumn losers: \n%v\n", r.VolumeLosers)
}

// GeneratePortfolio generate portfolio
func GeneratePortfolio(p *Portfolio) {

}
