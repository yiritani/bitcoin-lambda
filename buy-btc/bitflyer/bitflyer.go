package bitflyer

import (
	"encoding/hex"
	"encoding/json"
	"buy-btc/utils"

	"crypto/hmac"
	"crypto/sha256"
	"strconv"
	"time"
)

const baseURL = "https://api.bitflyer.com"
const productCodeKey = "product_code"

type Ticker struct {
	ProductCode			string `json:"product_code"`
	State 				string `json:"state"`
	Timestamp 			string `json:"timestamp"`
	TickID 				int	   `json:"tick_id"`
	BestBid 			float64 `json:"best_bid"`
	BestAsk 			float64 `json:"best_ask"`
	BestBidSize 		float64 `json:"best_bid_size"`
	BestAskSize 		float64 `json:"best_ask_size"`
	TotalBidDepth 		float64 `json:"total_bid_depth"`
	TotalAskDepth 		float64 `json:"total_ask_depth"`
	MarketBidSize 		float64 `json:"market_bid_size"`
	MarketAskSize 		float64 `json:"market_ask_size"`
	Ltp 				float64 `json:"ltp"`
	Volume 				float64 `json:"volume"`
	VolumeByProduct 	float64 `json:"volume_by_product"`
}


type Order struct {
	ProductCode 	string 	`json:"product_code"`
	ChildOrderType 	string 	`json:"child_order_type"`
	Side  			string 	`json:"side"`
	Price 			float64 `json:"price"`
	Size 			float64 `json:"size"`
	MinuteToExpire 	int 	`json:"minute_to_expire"`
	TimeInForce 	string 	`json:"time_in_force"`
}

type OrderRes struct {
	ChildOrderAcceptanceId string `json:"child_order_acceptance_id"`
}

func GetTicker(code ProductCode) (*Ticker, error) {
	url := baseURL + "/v1/ticker"
	res, err := utils.DoHttpRequest("GET", url, nil,
				map[string]string{productCodeKey: code.String()}, nil)
	if err != nil{
		return nil, err
	}

	var ticker Ticker
	err = json.Unmarshal(res, &ticker)
	if err != nil{
		return nil, err
	}

	return &ticker, nil
}

func getHeader(method, path, apiKey, apiSecret string, body []byte) map[string]string{
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	//ACCESS-SIGN は、ACCESS-TIMESTAMP, HTTP メソッド, リクエストのパス, リクエストボディを文字列として連結したものを、
	//API secret で HMAC-SHA256 署名を行った結果です。
	text := timestamp + method + path + string(body)
	mac := hmac.New(sha256.New, []byte(apiSecret))
	mac.Write([]byte(text))
	sign := hex.EncodeToString(mac.Sum(nil))

	return map[string]string{
		"ACCESS-KEY":			apiKey,
		"ACCESS-TIMESTAMP":		timestamp,
		"ACCESS-SIGN":			sign,
		"Content-Type":			"application/json",
	}
}
