package coinmarket

import "github.com/df-mc/atomic"

type Currency struct {
	price atomic.Value[float64]
}

func (c Currency) Price() float64 {
	return c.price.Load()
}

func BTC() Currency {
	return currency("BTC")
}
func LTC() Currency {
	return currency("LTC")
}
func DOGE() Currency {
	return currency("DOGE")
}
func ETH() Currency {
	return currency("ETH")
}

func currency(curr string) Currency {
	data.currencyDataMu.Lock()
	defer data.currencyDataMu.Unlock()
	return data.CurrencyData[curr]
}
