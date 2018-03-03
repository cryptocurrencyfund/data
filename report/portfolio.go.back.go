package report

import (
	"time"
)

// CoinEntry CoinEntry
type CoinEntry struct {
	Names       string
	Symbol      string
	CostBasis   float64
	Quantity    float64
	lastUpdated time.Time
}

// CoinHolding Coin Holding obj
type CoinHolding struct {
	Name           string  `json:"name"`
	Symbol         string  `json:"symbol"`
	CostBasis      float64 `json:"costBasis"`
	CostBasisBtc   float64
	MarketPrice    float64
	MarketPriceBtc float64
	Quantity       float64
	lastUpdated    time.Time
}

// TotalMarketValue TotalMarketValue
func (c *CoinHolding) TotalMarketValue() float64 {
	return c.MarketPrice * c.Quantity
}

// TotalMarketValueBtc TotalMarketValueBtc
func (c *CoinHolding) TotalMarketValueBtc() float64 {
	return c.MarketPriceBtc * c.Quantity
}

// TotalGainLoss TotalGainLoss
func (c *CoinHolding) TotalGainLoss() float64 {
	return (c.MarketPrice - c.CostBasis) * c.Quantity
}

// TotalGainLossBtc TotalGainLossBtc
func (c *CoinHolding) TotalGainLossBtc() float64 {
	return (c.MarketPriceBtc - c.CostBasisBtc) * c.Quantity
}

// Portfolio Portfolio obj
type Portfolio struct {
	Name           string
	Assests        []*CoinHolding
	Principle      float64
	ManagementFee  float64
	Expense        float64
	lastUpdated    time.Time
	lastRebalanced time.Time
}

// TotalMarketValue TotalMarketValue
func (p *Portfolio) TotalMarketValue() (sum float64) {
	sum = 0
	for i := range p.Assests {
		sum += p.Assests[i].TotalMarketValue()
	}
	return
}

// TotalMarketValueBtc TotalMarketValueBtc
func (p *Portfolio) TotalMarketValueBtc() (sum float64) {
	sum = 0
	for i := range p.Assests {
		sum += p.Assests[i].TotalMarketValueBtc()
	}
	return
}

// TotalGainLoss TotalGainLoss
func (p *Portfolio) TotalGainLoss() float64 {
	return p.TotalMarketValue() - p.Principle
}

// Rebalance Rebalance
func (p *Portfolio) Rebalance(c *CoinEntry) {

}

// Invest Invest
func (p *Portfolio) Invest(c *CoinEntry) {

}
