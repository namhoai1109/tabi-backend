package main

import (
	"fmt"
	"tabi-payment/internal/functions/migration"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(func() (string, error) {
		err := migration.Run()
		if err != nil {
			return "ERROR", fmt.Errorf("ERROR: %+v", err)
		}

		return "Migrated!", nil
	})
}
