package main

import (
	"context"
	"fmt"
	seedgeneraltype "tabi-booking/internal/functions/seed/generaltype"

	"github.com/namhoai1109/tabi/core/logger"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(func(ctx context.Context) (string, error) {
		err := seedgeneraltype.Run(ctx)
		if err != nil {
			logger.LogError(ctx, fmt.Sprintf("function seed general type run with err: %v", err))
			return "ERROR", fmt.Errorf("ERROR: %+v", err)
		}

		logger.LogInfo(ctx, "function seed general type run OK")
		return "function seed general type run OK", nil
	})
}
