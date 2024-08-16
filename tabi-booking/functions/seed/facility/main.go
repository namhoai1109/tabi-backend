package main

import (
	"context"
	"fmt"
	seedfacility "tabi-booking/internal/functions/seed/facility"

	"github.com/namhoai1109/tabi/core/logger"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(func(ctx context.Context) (string, error) {
		err := seedfacility.Run(ctx)
		if err != nil {
			logger.LogError(ctx, fmt.Sprintf("function seed facility run with err: %v", err))
			return "ERROR", fmt.Errorf("ERROR: %+v", err)
		}

		logger.LogInfo(ctx, "function seed facility run OK")
		return "function seed facility run OK", nil
	})
}
