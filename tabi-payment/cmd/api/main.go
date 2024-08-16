package main

import (
	"context"
	"fmt"
	"tabi-payment/config"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/labstack/echo/v4"

	"github.com/namhoai1109/tabi/core/middleware/jwt"
	"github.com/namhoai1109/tabi/core/paypal"
	"github.com/namhoai1109/tabi/core/server"

	userautho "tabi-payment/internal/api/user/autho"
	userpayment "tabi-payment/internal/api/user/payment"
	s2s "tabi-payment/internal/s2s"
)

var echoLambda *echoadapter.EchoLambda

func main() {
	fmt.Println("Start lambda function...")
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	fmt.Println("Loaded config!")

	// db, err := dbutil.New(cfg.DbDsn, cfg.DbLog)
	// if err != nil {
	// 	panic(err)
	// }

	// sqlDB, err := db.DB()
	// if err != nil {
	// 	panic(err)
	// }
	// defer sqlDB.Close()
	// fmt.Println("Connected to DB!")

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
	// db.Logger = logadapter.NewGormLogger().LogMode(gormlogger.Info)

	// * swagger ui
	e.Static("/swaggerui", "swaggerui")

	// * ping route
	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(200, "server is running!")
	})

	// ====================== API ======================

	// * Authorizations
	userAuthoService := userautho.New()

	// * Initialize services
	// === jwt service ===
	jwtPartnerService := jwt.New(cfg.JwtPartnerAlgorithm, cfg.JwtPartnerSecret, cfg.JwtPartnerDuration)
	jwtUserService := jwt.New(cfg.JwtUserAlgorithm, cfg.JwtUserSecret, cfg.JwtUserDuration)

	// === general service ===
	paypalService := paypal.New(cfg.PaypalBaseURL, paypal.Config{
		ClientID:     cfg.PaypalClientID,
		ClientSecret: cfg.PaypalClientSecret,
		Debug:        cfg.ReqLog,
		Timeout:      20,
	})
	s2sService := s2s.New(cfg, jwtPartnerService)

	// === user services ===
	userPaymentService := userpayment.New(paypalService, s2sService)

	// * Initialize HTTP handlers
	// === user ===
	userRouter := e.Group("/user")
	userRouter.Use(jwtUserService.MiddlewareFunction())
	userpayment.NewHTTP(userPaymentService, userAuthoService, userRouter.Group("/payments"))

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
