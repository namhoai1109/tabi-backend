package partner

import (
	"tabi-booking/config"
	"tabi-booking/internal/usecase/branch"
	"time"

	dbcore "github.com/namhoai1109/tabi/core/db"

	"gorm.io/gorm"
)

// New creates new authen partner service
func New(
	db *gorm.DB,
	accountDB AccountDB,
	representativeDB RepresentativeDB,
	branchManagerDB BranchManagerDB,
	companyDB CompanyDB,
	branchDB BranchDB,
	jwt JWT,
	cfg *config.Configuration,
	branchUseCase branch.BranchUseCase,
) *AuthenPartner {
	return &AuthenPartner{
		db:               db,
		accountDB:        accountDB,
		representativeDB: representativeDB,
		branchManagerDB:  branchManagerDB,
		companyDB:        companyDB,
		branchDB:         branchDB,
		jwt:              jwt,
		cfg:              cfg,
		branchUseCase:    branchUseCase,
	}
}

type AuthenPartner struct {
	db               *gorm.DB
	accountDB        AccountDB
	representativeDB RepresentativeDB
	branchManagerDB  BranchManagerDB
	companyDB        CompanyDB
	branchDB         BranchDB
	jwt              JWT
	cfg              *config.Configuration
	branchUseCase    branch.BranchUseCase
}

type AccountDB interface {
	dbcore.Intf
}

type RepresentativeDB interface {
	dbcore.Intf
}

type BranchManagerDB interface {
	dbcore.Intf
}

type CompanyDB interface {
	dbcore.Intf
}

type JWT interface {
	GenerateToken(claims map[string]interface{}, expire *time.Time) (string, int, error)
}

type BranchDB interface {
	dbcore.Intf
}
