package main

import (
	"context"
	"fmt"

	"tabi-file/config"
	partnerautho "tabi-file/internal/api/partner/autho"
	partnerfile "tabi-file/internal/api/partner/file"
	publicfile "tabi-file/internal/api/public/file"
	filedb "tabi-file/internal/db/file"
	"tabi-file/internal/s3"
	dbutil "tabi-file/internal/util/db"

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

	// initialize db instance
	fileDB := filedb.NewDB()

	// initialize general services
	jwtPartnerSvc := jwt.New(cfg.JWTPartnerAlgorithm, cfg.JWTPartnerSecret, cfg.JWTPartnerDuration)
	// Initialize S3 services
	s3Svc := s3.New(s3.Config{
		Region:          cfg.S3Region,
		AccessKeyID:     cfg.S3AccessKeyID,
		SecretAccessKey: cfg.S3SecretAccessKey,
		BucketName:      cfg.S3PublicBucketName,
		Debug:           true,
	})
	fmt.Println("Initialized S3 service!")

	// initialize partner services
	partnerAuthoSvc := partnerautho.New()
	partnerFileSvc := partnerfile.New(db, fileDB, s3Svc, cfg)
	publicFileSvc := publicfile.New(db, fileDB, s3Svc, cfg)

	// initialize partner router
	partnerRouter := e.Group("/partner")
	partnerRouter.Use(jwtPartnerSvc.MiddlewareFunction())
	partnerfile.NewHTTP(partnerFileSvc, partnerAuthoSvc, partnerRouter.Group("/files"))

	// initialize public router
	publicfile.NewHTTP(publicFileSvc, e.Group("/files"))

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
