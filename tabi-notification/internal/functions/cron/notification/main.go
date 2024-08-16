package notification

import (
	"context"
	"fmt"
	"tabi-notification/config"
	"tabi-notification/internal/expo"
	"tabi-notification/internal/model"
	"tabi-notification/internal/s2s"
	"time"

	"github.com/namhoai1109/tabi/core/logger"
)

func Run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	fmt.Println("Loaded config!")

	expoSvc := expo.New()

	s2s := s2s.New(cfg)

	now := time.Now()
	date := now.Format("2006-01-02")
	hour := now.Hour() + 1
	notificationTime := fmt.Sprintf("%sT%02d:00:00Z", date, hour)

	schedules, err := s2s.GetScheduleList(ctx, map[string]interface{}{
		"time":        notificationTime,
		"is_notified": false,
	})
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("error getting schedule list: %v", err))
		return err
	}

	userIDs := []int{}
	destinationMap := map[int]*model.Schedule{}
	for _, schedule := range schedules {
		userIDs = append(userIDs, schedule.UserID)
		destinationMap[schedule.UserID] = schedule
	}

	devices, err := s2s.GetDeviceList(ctx, map[string]interface{}{
		"user_id__in": userIDs,
		"is_active":   true,
	})
	if err != nil {
		logger.LogError(ctx, fmt.Sprintf("error getting device list: %v", err))
		return err
	}

	logger.LogInfo(ctx, fmt.Sprintf("Found %d devices", len(devices)))

	for _, device := range devices {

		if destinationMap[device.UserID] == nil {
			continue
		}

		title := fmt.Sprintf("%s - %s", destinationMap[device.UserID].DestinationName, destinationMap[device.UserID].DestinationCategory)

		time := time.Time(destinationMap[device.UserID].StartTime)
		body := fmt.Sprintf("You have a schedule at %s at %d:%d. Don't forget to check it out!", destinationMap[device.UserID].DestinationName, time.Hour(), time.Minute())

		if err := expoSvc.SendNotification(device.PushToken, title, body); err != nil {
			logger.LogError(ctx, fmt.Sprintf("error sending push notification to device %s: %v", device.PushToken, err))
		}

		if err := s2s.MarkNotified(ctx, destinationMap[device.UserID].ID); err != nil {
			logger.LogError(ctx, fmt.Sprintf("error marking schedule as notified: %v", err))
		}
	}

	return nil
}
