package bitflyer

import (
	"buy-btc/utils"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

const baseURL = "https://api.bitflyer.com"
const productCodeKey = "product_code"
const btcMinimumAmount = 0.001 //bitflyerの最小注文数量
const btcPlace = 4.0           //少数第4位

type APIClient struct {
	apiKey    string
	apiSecret string
}

func NewAPIClient(apiKey, apiSecret string) *APIClient {
	return &APIClient{apiKey, apiSecret}
}

type Ticker struct {
	ProductCode     string  `json:"product_code"`
	State           string  `json:"state"`
	Timestamp       string  `json:"timestamp"`
	TickID          int     `json:"tick_id"`
	BestBid         float64 `json:"best_bid"`
	BestAsk         float64 `json:"best_ask"`
	BestBidSize     float64 `json:"best_bid_size"`
	BestAskSize     float64 `json:"best_ask_size"`
	TotalBidDepth   float64 `json:"total_bid_depth"`
	TotalAskDepth   float64 `json:"total_ask_depth"`
	MarketBidSize   float64 `json:"market_bid_size"`
	MarketAskSize   float64 `json:"market_ask_size"`
	Ltp             float64 `json:"ltp"`
	Volume          float64 `json:"volume"`
	VolumeByProduct float64 `json:"volume_by_product"`
}

type Order struct {
	ProductCode    string  `json:"product_code"`
	ChildOrderType string  `json:"child_order_type"`
	Side           string  `json:"side"`
	Price          float64 `json:"price"`
	Size           float64 `json:"size"`
	MinuteToExpire int     `json:"minute_to_expire"`
	TimeInForce    string  `json:"time_in_force"`
}

type OrderRes struct {
	ChildOrderAcceptanceId string `json:"child_order_acceptance_id"`
}

func GetTicker(ch chan *Ticker, errCh chan error, code ProductCode) {
	url := baseURL + "/v1/ticker"
	res, err := utils.DoHttpRequest("GET", url, nil,
		map[string]string{productCodeKey: code.String()}, nil)
	if err != nil {
		errCh <- err
		ch <- nil
		return
	}

	var ticker Ticker
	err = json.Unmarshal(res, &ticker)
	if err != nil {
		errCh <- err
		ch <- nil
		return
	}

	ch <- &ticker
	errCh <- nil

}

func PlaceOrderWithParams(client *APIClient, price, size float64) (*OrderRes, error) {
	order := Order{
		ProductCode:    Btcjpy.String(),
		ChildOrderType: Limit.String(),
		Side:           Buy.String(),
		Price:          price,
		Size:           size,
		MinuteToExpire: 4320,
		TimeInForce:    Gtc.String(),
	}

	orderRes, err := client.PlaceOrder(&order)
	if err != nil {
		return nil, err
	}

	return orderRes, nil
}

func GetBuyLogic(strategy int) func(float64, *Ticker) (float64, float64) {
	var logic func(float64, *Ticker) (float64, float64)

	switch strategy {
	case 1:
		logic = func(budget float64, t *Ticker) (float64, float64) {
			var buyPrice, buySize float64
			buyPrice = utils.RoundDecimal(t.Ltp * 0.985)
			buySize = utils.CalcAmount(buyPrice, budget, btcMinimumAmount, btcPlace)
			return buyPrice, buySize
		}
		break
	default:
		logic = func(budget float64, t *Ticker) (float64, float64) {
			var buyPrice, buySize float64
			buyPrice = utils.RoundDecimal(t.BestAsk)
			buySize = utils.CalcAmount(buyPrice, budget, btcMinimumAmount, btcPlace)
			return buyPrice, buySize
		}
		break
	}

	return logic
}

func (client *APIClient) PlaceOrder(order *Order) (*OrderRes, error) {
	method := "POST"
	path := "/v1/me/sendchildorder"
	url := baseURL + path
	data, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}

	header := client.getHeader(method, path, data)

	res, err := utils.DoHttpRequest(method, url, header, map[string]string{}, data)
	if err != nil {
		return nil, err
	}

	var orderRes OrderRes
	err = json.Unmarshal(res, &orderRes)
	if err != nil {
		return nil, err
	}

	if len(orderRes.ChildOrderAcceptanceId) == 0 {
		return nil, errors.New(string(res))
	}

	return &orderRes, nil
}

func (client *APIClient) getHeader(method, path string, body []byte) map[string]string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	//ACCESS-SIGN は、ACCESS-TIMESTAMP, HTTP メソッド, リクエストのパス, リクエストボディを文字列として連結したものを、
	//API secret で HMAC-SHA256 署名を行った結果です。
	text := timestamp + method + path + string(body)
	mac := hmac.New(sha256.New, []byte(client.apiSecret))
	mac.Write([]byte(text))
	sign := hex.EncodeToString(mac.Sum(nil))

	return map[string]string{
		"ACCESS-KEY":       client.apiKey,
		"ACCESS-TIMESTAMP": timestamp,
		"ACCESS-SIGN":      sign,
		"Content-Type":     "application/json",
	}
}
