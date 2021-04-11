package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	lmb "github.com/suzuki-shunsuke/aws-iam-cred-sender/pkg/lambda"
)

func main() {
	handler := lmb.Handler{}
	ctx := context.Background()
	if err := handler.Init(ctx); err != nil {
		log.Fatal(err)
	}
	lambda.Start(handler.Start)
}
