package coinmarket

import (
	"github.com/df-mc/atomic"
	"github.com/pragu3/gophig"
	"log"
	"os"
	"sync"
)

type currencyData struct {
	Price float64 `json:"price"`
}

func (c currencyData) currency() Currency {
	return Currency{
		price: *atomic.NewValue(c.Price),
	}
}

func (c Currency) data() currencyData {
	return currencyData{
		Price: c.Price(),
	}
}

type marketData struct {
	currencyDataMu sync.Mutex
	CurrencyData   map[string]Currency `json:"currency_data"`
}

func newMarketData() *marketData {
	m := &marketData{
		CurrencyData: map[string]Currency{},
	}

	m.currencyDataMu.Lock()
	var dat = map[string]currencyData{}
	err := conf.GetConf(&dat)
	if err != nil {
		if os.IsNotExist(err) {
			if err = m.save(); err != nil {
				log.Fatalln(err)
			}
			return nil
		}
		log.Fatalln(err)
		return nil
	}
	for s, c := range dat {
		m.CurrencyData[s] = c.currency()
	}
	m.currencyDataMu.Unlock()
	return m
}

func (m *marketData) save() error {
	var d = map[string]currencyData{}
	for s, c := range m.CurrencyData {
		d[s] = c.data()
	}
	return conf.SetConf(d)
}

var (
	conf = gophig.NewGophig("assets/currency_data", "json", 740)
	data = newMarketData()
)
