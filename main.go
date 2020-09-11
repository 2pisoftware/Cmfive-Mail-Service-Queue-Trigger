package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/2pisoftware/mail-service-popper/handler"
)

func main() {
	lambda.Start(handler.Handle)
}
