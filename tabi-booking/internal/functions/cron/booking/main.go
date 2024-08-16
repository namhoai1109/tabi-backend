package booking

import (
	"context"
	"fmt"
	"tabi-booking/config"
	bookingdb "tabi-booking/internal/db/booking"
	"tabi-booking/internal/model"
	dbutil "tabi-booking/internal/util/db"
	"time"

	"github.com/imdatngo/gowhere"
	dbcore "github.com/namhoai1109/tabi/core/db"
	"github.com/namhoai1109/tabi/core/logger"
	"github.com/namhoai1109/tabi/core/middleware/logadapter"
	gormlogger "gorm.io/gorm/logger"
)

func Run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	fmt.Println("Loaded config!")

	db, err := dbutil.New(cfg.DbDsn, cfg.DbLog)
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()
	fmt.Println("Connected to DB!")
	// log adapter
	db.Logger = logadapter.NewGormLogger().LogMode(gormlogger.Info)

	// * Initialize DB interfaces
	bookingDB := bookingdb.NewDB()

	lq := &dbcore.ListQueryCondition{
		Filter: gowhere.WithConfig(gowhere.Config{Strict: true}),
	}
	lq.Filter.And(`status IN (?)`, []string{model.BookingStatusApproved, model.BookingStatusInReview})
	bookings := []*model.Booking{}
	if err := bookingDB.List(db, &bookings, lq, nil); err != nil {
		logger.LogError(ctx, fmt.Sprintf("Failed to get booking list: %s", err))
		return err
	}

	if len(bookings) == 0 {
		return nil
	}

	inReviewBookingIDs := []int{}
	completedBookingIDs := []int{}
	for _, booking := range bookings {
		fourDays := time.Hour * 24 * 4
		if booking.Status == model.BookingStatusInReview && time.Now().After((*booking.CheckOutDate).Add(fourDays)) {
			completedBookingIDs = append(completedBookingIDs, booking.ID)
			continue
		}

		if booking.Status == model.BookingStatusApproved && time.Now().After(*booking.CheckOutDate) {
			inReviewBookingIDs = append(inReviewBookingIDs, booking.ID)
			continue
		}
	}

	if len(inReviewBookingIDs) > 0 {
		if err := bookingDB.Update(db, map[string]interface{}{
			"status": model.BookingStatusInReview,
		}, `id IN (?)`, inReviewBookingIDs); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Failed to update booking status %s: %v", model.BookingStatusInReview, err))
		}
	}

	if len(completedBookingIDs) > 0 {
		if err := bookingDB.Update(db, map[string]interface{}{
			"status": model.BookingStatusCompleted,
		}, `id IN (?)`, completedBookingIDs); err != nil {
			logger.LogError(ctx, fmt.Sprintf("Failed to update booking status %s: %v", model.BookingStatusCompleted, err))
		}
	}

	return nil
}
