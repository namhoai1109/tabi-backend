package branch

import (
	"tabi-booking/internal/model"

	dbcore "github.com/namhoai1109/tabi/core/db"
	"gorm.io/gorm"
)

func NewDB() *DB {
	return &DB{dbcore.NewDB(&model.Branch{})}
}

type DB struct {
	*dbcore.DB
}

func (d *DB) CheckExitedBranchManager(db *gorm.DB, branchID int) (bool, error) {
	branch := &model.Branch{}
	if err := d.View(db.Preload("BranchManager"), &branch, `id = ?`, branchID); err != nil {
		return false, err
	}
	if branch.BranchManager != nil {
		return true, nil
	}

	return false, nil
}

func (d *DB) GetBranchIDByBranchManagerID(db *gorm.DB, managerID int) (int, error) {
	branch := &model.Branch{}
	if err := d.View(db, &branch, `branch_manager_id = ?`, managerID); err != nil {
		return 0, err
	}

	return branch.ID, nil
}

func (d *DB) GetBranchIDsByRepresentativeID(db *gorm.DB, companyID int) ([]int, error) {
	company := &model.Company{}
	branchIDs := []int{}
	if err := d.View(db.Preload("Branches"), &company, `representative_id = ?`, companyID); err != nil {
		return branchIDs, err
	}

	for _, branch := range company.Branches {
		branchIDs = append(branchIDs, branch.ID)
	}

	return branchIDs, nil
}

func (d *DB) AnalyzeBranchesRevenue(db *gorm.DB, branchIDs []int, year int) ([]*model.RevenueAnalysisData, error) {
	res := []*model.RevenueAnalysisData{}
	err := db.Raw(`
		SELECT DATE_PART('MONTH', b.check_in_date) AS month, SUM(b.total_price) AS revenue
		FROM public.booking b 
		WHERE deleted_at IS NULL 
		AND status IN (?)
		AND room_id IN (
			SELECT r.id 
			FROM public.room r 
			WHERE branch_id IN (?)
			AND deleted_at IS NULL
		)
		AND DATE_PART('YEAR', b.check_in_date) = ?
		GROUP BY DATE_PART('MONTH',b.check_in_date), b.total_price
		ORDER BY DATE_PART('MONTH',b.check_in_date) ASC
	`, []string{model.BookingStatusApproved}, branchIDs, year).Scan(&res).Error
	return res, err
}

func (d *DB) AnalyzeBookingRequestQuantity(db *gorm.DB, branchIDs []int, year int) ([]*model.BookingRequestQuantityAnalysisData, error) {
	res := []*model.BookingRequestQuantityAnalysisData{}
	err := db.Raw(`
		SELECT DATE_PART('MONTH',b.check_in_date) AS month, count(id) AS quantity
		FROM public.booking b 
		WHERE deleted_at IS NULL 
		AND status IN (?)
		AND room_id IN (
			SELECT r.id 
			FROM public.room r 
			WHERE branch_id IN (?)
			AND deleted_at IS NULL
		)
		AND DATE_PART('YEAR', b.check_in_date) = ?
		GROUP BY DATE_PART('MONTH',b.check_in_date)
		ORDER BY DATE_PART('MONTH',b.check_in_date) ASC
	`, []string{model.BookingStatusApproved}, branchIDs, year).Scan(&res).Error
	return res, err
}

func (d *DB) ListFeaturedDestination(db *gorm.DB) ([]string, error) {
	res := []string{}
	err := db.Raw(`
		SELECT b.province_city 
		FROM public.branch b 
		WHERE b.deleted_at IS NULL 
		AND b.id IN (
		SELECT r.branch_id AS branch_id
			FROM public.room r 
			JOIN public.booking b ON r.id = b.room_id AND b.deleted_at IS NULL
			WHERE r.deleted_at IS NULL 
			AND b.status IN (?)
			GROUP BY r.branch_id 
			ORDER BY count(b.id) DESC 
		)
		GROUP BY b.province_city 
	`, []string{model.BookingStatusApproved}).Scan(&res).Error
	return res, err
}

func (d *DB) ListFeaturedBranches(db *gorm.DB) ([]*model.Branch, error) {
	res := []*model.Branch{}
	err := db.Raw(`
		SELECT b.* 
		FROM public.branch b 
		WHERE b.deleted_at IS NULL 
		AND b.id IN (
		SELECT r.branch_id AS branch_id
			FROM public.room r 
			JOIN public.booking b ON r.id = b.room_id AND b.deleted_at IS NULL
			WHERE r.deleted_at IS NULL 
			AND b.status IN (?)
			GROUP BY r.branch_id 
			ORDER BY count(b.id) DESC 
		)
		LIMIT 10
	`, []string{model.BookingStatusApproved}).Scan(&res).Error
	return res, err
}
