package main

import (
	"context"
	"fmt"
	"tabi-notification/internal/functions/cron/notification"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/namhoai1109/tabi/core/logger"
)

func main() {
	lambda.Start(func(ctx context.Context) error {
		if err := notification.Run(ctx); err != nil {
			logger.LogError(ctx, fmt.Sprintf("error running notification cron job: %v", err))
			return err
		}

		logger.LogInfo(ctx, "notification cron job executed successfully")

		return nil
	})
}
