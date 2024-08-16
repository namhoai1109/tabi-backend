package main

import (
	"context"
	"fmt"
	"tabi-booking/internal/functions/cron/booking"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/namhoai1109/tabi/core/logger"
)

func main() {
	lambda.Start(func(ctx context.Context) error {
		if err := booking.Run(ctx); err != nil {
			logger.LogError(ctx, fmt.Sprintf("error running booking cron job: %v", err))
			return err
		}

		logger.LogInfo(ctx, "booking cron job executed successfully")

		return nil
	})
}

// run local
// func main() {
// 	ctx := context.Background()
// 	if err := booking.Run(ctx); err != nil {
// 		logger.LogError(ctx, fmt.Sprintf("error running booking cron job: %v", err))
// 		return
// 	}

// 	logger.LogInfo(ctx, "booking cron job executed successfully")
// }
