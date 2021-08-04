package main

import (
	//"buy-btc/bitflyer"
	"fmt"

	// "github.com/aws/aws-sdk-go"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"

)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//ticker, err := bitflyer.GetTicker(bitflyer.Btcjpy)

	apikey, err := getParameter("buy-btc-apikey")
	if err != nil {
		return getErrorResponse(err.Error()), err
	}

	//secretkey, err := getParameter("\tbuy-btc-apisecret")
	//if err != nil {
	//	return getErrorResponse(err.Error()), err
	//}

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Bad Request!!",
			StatusCode: 400,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Apikey:%+v", apikey),
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