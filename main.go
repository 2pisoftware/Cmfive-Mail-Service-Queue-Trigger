package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/2pisoftware/Cmfive-Mail-Service-Queue-Trigger/handler"
)

func main() {
	lambda.Start(handler.Handle)
}
