package main

import (
	"github.com/2pisoftware/mail-service-popper/popper"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(popper.HandlePop)
}
