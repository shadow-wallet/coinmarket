package coinmarket

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type CoinMarket struct {
	apiKey string
}

func New(apiKey string) *CoinMarket {
	return &CoinMarket{
		apiKey: apiKey,
	}
}

func (c *CoinMarket) UpdateFiat() {
	cl := http.DefaultClient
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest?convert=CAD", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("X-CMC_PRO_API_KEY", c.apiKey)
	req.Header.Add("Accepts", "application/json")

	resp, err := cl.Do(req)
	if err != nil {
		panic(err)
	}

	var res map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(body, &res)

	data.currencyDataMu.Lock()
	for _, d := range res["data"].([]interface{}) {
		buf, err := json.Marshal(d.(map[string]interface{})["quote"].(map[string]interface{})["CAD"])
		if err != nil {
			log.Fatalln(err)
		}
		var currData currencyData
		err = json.Unmarshal(buf, &currData)
		if err != nil {
			log.Fatalln(err)
		}
		data.CurrencyData[d.(map[string]interface{})["symbol"].(string)] = currData.currency()
	}
	data.currencyDataMu.Unlock()

	err = data.save()
	if err != nil {
		log.Fatalln(err)
	}
}
