package main

import (
	"context"
	"fmt"

	"tabi-booking/config"
	authenpartner "tabi-booking/internal/api/authen/partner"
	authenuser "tabi-booking/internal/api/authen/user"
	partnerautho "tabi-booking/internal/api/partner/autho"
	partnerbank "tabi-booking/internal/api/partner/bank"
	partnerbooking "tabi-booking/internal/api/partner/booking"
	partnerbranch "tabi-booking/internal/api/partner/branch"
	branchmanager "tabi-booking/internal/api/partner/branchmanager"
	partnercompany "tabi-booking/internal/api/partner/company"
	partnerme "tabi-booking/internal/api/partner/me"
	partnerroom "tabi-booking/internal/api/partner/room"
	partnerroomtype "tabi-booking/internal/api/partner/roomtype"
	publicbranch "tabi-booking/internal/api/public/branch"
	publicfacility "tabi-booking/internal/api/public/facility"
	publicgeneraltype "tabi-booking/internal/api/public/generaltype"
	publicuser "tabi-booking/internal/api/public/user"
	userautho "tabi-booking/internal/api/user/autho"
	userbooking "tabi-booking/internal/api/user/booking"
	userbranch "tabi-booking/internal/api/user/branch"
	userme "tabi-booking/internal/api/user/me"
	usersurvey "tabi-booking/internal/api/user/survey"
	accountdb "tabi-booking/internal/db/account"
	bankdb "tabi-booking/internal/db/bank"
	bookingdb "tabi-booking/internal/db/booking"
	branchdb "tabi-booking/internal/db/branch"
	branchmanagerdb "tabi-booking/internal/db/branchmanager"
	companydb "tabi-booking/internal/db/company"
	facilitydb "tabi-booking/internal/db/facility"
	facturereductiondb "tabi-booking/internal/db/facturereduction"
	generaltypedb "tabi-booking/internal/db/generaltype"
	ratingdb "tabi-booking/internal/db/rating"
	representativedb "tabi-booking/internal/db/representative"
	reservationreductiondb "tabi-booking/internal/db/reservationreduction"
	roomdb "tabi-booking/internal/db/room"
	roomtypedb "tabi-booking/internal/db/roomtype"
	roomtypeofbranchdb "tabi-booking/internal/db/roomtypeofbranch"
	savedbranchdb "tabi-booking/internal/db/savedbranch"
	surveydb "tabi-booking/internal/db/survey"
	userdb "tabi-booking/internal/db/user"
	"tabi-booking/internal/rbac"
	dbutil "tabi-booking/internal/util/db"

	branchusecase "tabi-booking/internal/usecase/branch"

	"github.com/namhoai1109/tabi/core/middleware/jwt"
	"github.com/namhoai1109/tabi/core/middleware/logadapter"

	"github.com/namhoai1109/tabi/core/server"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/labstack/echo/v4"
	gormlogger "gorm.io/gorm/logger"
)

var echoLambda *echoadapter.EchoLambda

func main() {
	fmt.Println("Start lambda function...")
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

	// * Initialize HTTP server
	e := server.New(&server.Config{
		Stage:        cfg.Stage,
		Port:         cfg.Port,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		AllowOrigins: cfg.AllowOrigins,
	}, cfg.ReqLog)
	fmt.Println("Initialized HTTP server")

	// log adapter
	db.Logger = logadapter.NewGormLogger().LogMode(gormlogger.Info)

	// * swagger ui
	e.Static("/swaggerui", "swaggerui")

	// * ping route
	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(200, "server is running!")
	})

	// ====================== API ======================

	// * Initialize DB interfaces
	accountDB := accountdb.NewDB()
	representativeDB := representativedb.NewDB()
	companyDB := companydb.NewDB()
	branchDB := branchdb.NewDB()
	facilityDB := facilitydb.NewDB()
	bankDB := bankdb.NewDB()
	generalTypeDB := generaltypedb.NewDB()
	userDB := userdb.NewDB()
	roomTypeDB := roomtypedb.NewDB()
	roomTypeOfBranchDB := roomtypeofbranchdb.NewDB()
	branchManagerDB := branchmanagerdb.NewDB()
	roomDB := roomdb.NewDB()
	factureReductionDB := facturereductiondb.NewDB()
	reservationReductionDB := reservationreductiondb.NewDB()
	savedBranchDB := savedbranchdb.NewDB()
	bookingDB := bookingdb.NewDB()
	surveyDB := surveydb.NewDB()
	ratingDB := ratingdb.NewDB()

	// * Authorizations
	partnerAuthoService := partnerautho.New()
	userAuthoService := userautho.New()

	// * Initialize use cases
	branchUseCase := branchusecase.New(db, branchDB, roomDB, facilityDB, generalTypeDB, bankDB)

	// * Initialize services
	// === jwt service ===
	jwtPartnerService := jwt.New(cfg.JwtPartnerAlgorithm, cfg.JwtPartnerSecret, cfg.JwtPartnerDuration)
	jwtUserService := jwt.New(cfg.JwtUserAlgorithm, cfg.JwtUserSecret, cfg.JwtUserDuration)

	// === authen service ===
	authenPartnerService := authenpartner.New(db, accountDB, representativeDB, branchManagerDB, companyDB, branchDB, jwtPartnerService, cfg, branchUseCase)
	authenUserService := authenuser.New(db, accountDB, userDB, jwtUserService, cfg)
	rbacSvc := rbac.New(cfg.ReqLog)

	// === partner service ===
	partnerBranchService := partnerbranch.New(db, accountDB, branchManagerDB, companyDB, branchDB, bankDB, facilityDB, generalTypeDB, branchUseCase, rbacSvc)
	partnerCompanyService := partnercompany.New(db, accountDB, representativeDB, companyDB, branchDB, cfg)
	partnerMeService := partnerme.New(db, representativeDB, branchManagerDB, companyDB, branchDB)
	roomTypeService := partnerroomtype.New(db, roomTypeDB, roomTypeOfBranchDB, branchDB, companyDB, facilityDB, rbacSvc)
	branchManagerService := branchmanager.New(db, accountDB, branchManagerDB, branchDB, companyDB)
	partnerBankService := partnerbank.New(db, bankDB, companyDB, rbacSvc)
	partnerRoomService := partnerroom.New(db, roomDB, roomTypeOfBranchDB, branchDB, generalTypeDB, factureReductionDB, reservationReductionDB, companyDB, facilityDB, bookingDB, userDB, rbacSvc)
	partnerBookingService := partnerbooking.New(db, bookingDB, roomDB, userDB, branchDB, rbacSvc)

	// === user service ===
	userMeService := userme.New(db, userDB)
	userBranchService := userbranch.New(db, branchDB, savedBranchDB, branchUseCase, ratingDB, bookingDB)
	userSurveyService := usersurvey.New(db, surveyDB)
	userBookingService := userbooking.New(db, roomDB, bookingDB)

	// === public service ===
	publicBranchService := publicbranch.New(db, branchDB, roomDB, facilityDB, branchUseCase, bookingDB, ratingDB)
	publicUserService := publicuser.New(db, userDB, surveyDB)
	publicFacilityService := publicfacility.New(db, facilityDB)
	publicGeneralTypeService := publicgeneraltype.New(db, generalTypeDB)

	// * Initialize HTTP handlers
	// === auth ===
	authenRouter := e.Group("/authen")
	authenpartner.NewHTTP(authenPartnerService, authenRouter.Group("/partners"))
	authenuser.NewHTTP(authenUserService, authenRouter.Group("/users"))

	// === partner ===
	partnerRouter := e.Group("/partner")
	partnerRouter.Use(jwtPartnerService.MiddlewareFunction())
	partnerbranch.NewHTTP(partnerBranchService, partnerAuthoService, partnerRouter.Group("/branches"))
	partnercompany.NewHTTP(partnerCompanyService, partnerAuthoService, partnerRouter.Group("/company"))
	partnerme.NewHTTP(partnerMeService, partnerAuthoService, partnerRouter.Group("/me"))
	partnerroomtype.NewHTTP(roomTypeService, partnerAuthoService, partnerRouter.Group("/room-types"))
	partnerbank.NewHTTP(partnerBankService, partnerAuthoService, partnerRouter.Group("/banks"))
	branchmanager.NewHTTP(branchManagerService, partnerAuthoService, partnerRouter.Group("/branch-managers"))
	partnerroom.NewHTTP(partnerRoomService, partnerAuthoService, partnerRouter.Group("/rooms"))
	partnerbooking.NewHTTP(partnerBookingService, partnerAuthoService, partnerRouter.Group("/bookings"))

	// === user ===
	userRouter := e.Group("/user")
	userRouter.Use(jwtUserService.MiddlewareFunction())
	userme.NewHTTP(userMeService, userAuthoService, userRouter.Group("/me"))
	userbranch.NewHTTP(userBranchService, userAuthoService, userRouter.Group("/branches"))
	usersurvey.NewHTTP(userSurveyService, userAuthoService, userRouter.Group("/surveys"))
	userbooking.NewHTTP(userBookingService, userAuthoService, userRouter.Group("/bookings"))

	// === public ===
	publicbranch.NewHTTP(publicBranchService, e.Group("/branches"))
	publicuser.NewHTTP(publicUserService, e.Group("/users"))
	publicfacility.NewHTTP(publicFacilityService, e.Group("/facilities"))
	publicgeneraltype.NewHTTP(publicGeneralTypeService, e.Group("/general-types"))

	// Start the HTTP server
	if cfg.Stage == "development" {
		// Start the HTTP server
		fmt.Println("Starting HTTP server...")
		server.Start(e, cfg.Stage == "development")
	} else {
		fmt.Println("Starting echoLambda...")
		echoLambda = echoadapter.New(e)
		lambda.Start(Handler)
	}
}

// Handler function
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return echoLambda.ProxyWithContext(ctx, req)
}
