package main

import (
	"fmt"
	"log"

	"github.com/career-built/marketPlace/Api"
	"github.com/career-built/marketPlace/Database"
	"github.com/career-built/marketPlace/MessageBroker"
	"github.com/career-built/marketPlace/OpenTelemetry"
	"github.com/career-built/marketPlace/Product_Service"
	"github.com/labstack/echo/v4"
)

func main() {
	fmt.Println("Starting Base")
	// define the type of dbConnector
	dbConnector := Database.NewPostgres()
	if dbConnector == nil {
		log.Fatal("can't connect to database")
	}
	defer dbConnector.CloseDB()

	// define the type of Message Broker
	messagebroker, err := MessageBroker.NewRabbitMQBroker("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("can't connect to messagebroker:: %s", err)
		panic(err)
	}
	defer messagebroker.Close()

	//inject the Database to the base Db interface
	productMgr := Product_Service.NewProductService(dbConnector, messagebroker)

	// define the type of Message Broker tracer
	messageTracer, err := OpenTelemetry.NewGlobalJaegerTracer("Test-tracing-app")
	if err != nil {
		log.Fatalf("can't connect to messageTracer:: %s", err)
		panic(err)
	}
	defer messageTracer.Close()
	// define the type of JWT
	jwtobj := JWT.newkeyCloakMiddleware()
	//inject the product feature to the base manager interface
	productRouter := Api.NewProductRouter(productMgr, messageTracer)
	e := echo.New()

	e.POST("/product/create", productRouter.CreateProduct)
	e.GET("/product/:id", productRouter.GetProductByID)
	e.GET("/getproductfromMarket", productRouter.GetProductFromMarket)

	e.Logger.Fatal(e.Start(":3030"))
}
