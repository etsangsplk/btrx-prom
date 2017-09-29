package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Example of a single market summary:
// {
// 			"MarketName" : "BTC-888",
// 			"High" : 0.00000919,
// 			"Low" : 0.00000820,
// 			"Volume" : 74339.61396015,
// 			"Last" : 0.00000820,
// 			"BaseVolume" : 0.64966963,
// 			"TimeStamp" : "2014-07-09T07:19:30.15",
// 			"Bid" : 0.00000820,
// 			"Ask" : 0.00000831,
// 			"OpenBuyOrders" : 15,
// 			"OpenSellOrders" : 15,
// 			"PrevDay" : 0.00000821,
// 			"Created" : "2014-03-20T06:00:00",
// 			"DisplayMarketName" : null
// }
type marketSummary struct {
	MarketName string
	Last       float64
	Volume     float64
}

func (m marketSummary) String() string {
	return fmt.Sprintf(`last{exchange="bittrex",market="%s"} %f`, m.MarketName, m.Last)
}

type marketSummaries struct {
	Success bool
	Message string
	Result  []marketSummary
}

func handler(w http.ResponseWriter, r *http.Request) {
	url := "https://bittrex.com/api/v1.1/public/getmarketsummaries"
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching market data: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unexpected error reading market data: %v", err), http.StatusInternalServerError)
		return
	}

	var summaries marketSummaries
	err = json.Unmarshal(data, &summaries)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unexpected error parsing market data: %v", err), http.StatusInternalServerError)
		return
	}

	if summaries.Success != true {
		fmt.Println(string(data))
		http.Error(w, fmt.Sprintf("Unsuccessful response from bittrex: %s", summaries.Message), http.StatusInternalServerError)
		return
	}

	for _, s := range summaries.Result {
		fmt.Fprintln(w, s)
	}

	fmt.Println("Data fetched successfully")
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe("127.0.0.1:8080", nil)
}
