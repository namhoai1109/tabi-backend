package main

import (
	"context"
	"fmt"
	"tabi-notification/config"
	publicdevice "tabi-notification/internal/api/device"
	publicschedule "tabi-notification/internal/api/schedule"
	dbutil "tabi-notification/internal/util/db"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/labstack/echo/v4"
	"github.com/namhoai1109/tabi/core/middleware/jwt"
	"github.com/namhoai1109/tabi/core/middleware/logadapter"
	"github.com/namhoai1109/tabi/core/server"
	gormlogger "gorm.io/gorm/logger"

	userautho "tabi-notification/internal/api/user/autho"
	userdevice "tabi-notification/internal/api/user/device"
	userschedule "tabi-notification/internal/api/user/schedule"

	devicedb "tabi-notification/internal/db/device"
	scheduledb "tabi-notification/internal/db/schedule"
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
	deviceDB := devicedb.NewDB()
	scheduleDB := scheduledb.NewDB()

	// * Authorizations
	userAuthoService := userautho.New()

	// * Initialize services
	// === jwt service ===
	jwtUserService := jwt.New(cfg.JwtUserAlgorithm, cfg.JwtUserSecret, cfg.JwtUserDuration)

	// === user service ===
	userScheduleSvc := userschedule.New(db, scheduleDB, cfg)
	userDeviceSvc := userdevice.New(db, deviceDB)

	// === public service ===
	deviceSvc := publicdevice.New(db, deviceDB)
	scheduleSvc := publicschedule.New(db, scheduleDB)

	// * Initialize HTTP handlers
	// === user ===
	userRouter := e.Group("/user")
	userRouter.Use(jwtUserService.MiddlewareFunction())
	userschedule.NewHTTP(userScheduleSvc, userAuthoService, userRouter.Group("/schedules"))
	userdevice.NewHTTP(userDeviceSvc, userAuthoService, userRouter.Group("/device"))

	// === public ===
	publicdevice.NewHTTP(deviceSvc, e.Group("/devices"))
	publicschedule.NewHTTP(scheduleSvc, e.Group("/schedules"))

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
