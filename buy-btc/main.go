package main

import (
	"buy-btc/bitflyer"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	tickerChan := make(chan *bitflyer.Ticker)
	errChan := make(chan error)
	defer close(tickerChan)
	defer close(errChan)

	go bitflyer.GetTicker(tickerChan, errChan, bitflyer.Btcjpy)
	ticker := <-tickerChan
	err := <-errChan

	if err != nil {
		return getErrorResponse(err.Error()), nil
	}

	apikey, err := getParameter("buy-btc-apikey")
	if err != nil {
		return getErrorResponse(err.Error()), err
	}

	apiSecret, err := getParameter("buy-btc-apisecret")
	if err != nil {
		return getErrorResponse(err.Error()), err
	}

	client := bitflyer.NewAPIClient(apikey, apiSecret)

	//全ての注文を削除する
	cancelRes, err := bitflyer.CancelOrderWithParams(client)
	fmt.Println("Cancel all orders", cancelRes)
	if err != nil {
		return getErrorResponse(err.Error()), err
	}

	//カリー化
	price, size := bitflyer.GetBuyLogic(1)(5000.0, ticker)
	orderRes, err := bitflyer.PlaceOrderWithParams(client, price, size)

	if err != nil {
		return getErrorResponse(err.Error()), err
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("res:%+v", orderRes),
		StatusCode: 200,
	}, nil
}

//Systems Managerからパラメータを取得する関数
func getParameter(key string) (string, error) {

	// ローカルのaws configを取得する
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := ssm.New(sess, aws.NewConfig().WithRegion("us-west-2"))

	params := &ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(true),
	}

	res, err := svc.GetParameter(params)
	if err != nil {
		return "", err
	}

	return *res.Parameter.Value, nil
}

func getErrorResponse(message string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       message,
		StatusCode: 400,
	}
}

func main() {
	lambda.Start(handler)
}
