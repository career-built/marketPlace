package api

import (
	"context"
	"example/baseProject/product"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// Instrumentation Based Tracing Interface
type MessageTracer interface {
	Close()
	StartGlobalSpan(operationtrachingNAme string) opentracing.Span
	StartConfig(serviceName string) (opentracing.Tracer, io.Closer, error)
}
type ProductManager interface {
	Add(product *product.Product, queueName string) error
	GetByID(id int) (*product.Product, error)
}

type ProductRouter struct {
	productManager ProductManager
	messageTracer  MessageTracer
}

func NewProductRouter(productManager ProductManager, messageTracer MessageTracer) *ProductRouter {
	return &ProductRouter{
		productManager: productManager,
		messageTracer:  messageTracer,
	}
}

func (obj *ProductRouter) CreateProduct(c echo.Context) error {
	//Request handeling
	product := new(product.Product)

	if err := c.Bind(product); err != nil {
		fmt.Printf("Error While Binding the product\n")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	//Service invoc
	if err := obj.productManager.Add(product, ""); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Respond handeling
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Product created successfully",
		"name":    product.NAME,
	})
}
func (obj *ProductRouter) GetProductByID(c echo.Context) error {

	productIDStr := c.Param("id")
	// Check if the productID is a valid integer
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid product ID",
		})
	}
	// Fetch the corresponding product from the database
	fmt.Printf("new Requested product ID: %d\n", productID)
	product, err := obj.productManager.GetByID(productID)
	if err == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Unaple to Fetch product from DB",
		})
	}
	if product == nil {
		return c.JSON(http.StatusOK, map[string]string{
			"INFO": "Product Not Found",
		})
	}
	// Respond with a JSON message
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Product fetched successfully",
		"id":      strconv.Itoa(product.ID),
		"name":    product.NAME,
		"price":   strconv.Itoa(product.PRICE),
	})
}

func (obj *ProductRouter) GetProductFromMarket(c echo.Context) error {
	//Configer jaeger
	fmt.Println("Start Config")
	serviceName := "product maneger servicen"
	mainspanString := "Get Product From Market Span"
	tracer, closer, err := obj.messageTracer.StartConfig(serviceName)
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
	// Start Span
	spanCtx, _ := tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(c.Request().Header),
	)
	span := tracer.StartSpan(mainspanString, ext.RPCServerOption(spanCtx))
	defer span.Finish()
	ctx := opentracing.ContextWithSpan(context.Background(), span)

	// //_____ Working Path ______
	categories := cattegories_api(ctx, mainspanString)
	// Loob throw categories and retrieve the products.
	for _, category := range categories {
		fmt.Println("-----------------loob throw category---------------------------------: ", category)
		// product_api(ctx, mainspanString, category)
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8081/marketservice", nil)

		// Inject the current span context into the outgoing HTTP headers
		span.Tracer().Inject(
			span.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(req.Header),
		)
		// Perform the HTTP request to service-b
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Error calling service-b:", err)
			return c.String(http.StatusInternalServerError, "Error calling service-b")
		}
		defer resp.Body.Close()
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "GetProductFromMarket successfully",
	})
}

func cattegories_api(ctx context.Context, mainspanString string) []string {
	fmt.Println("Getting Categories")
	span, _ := opentracing.StartSpanFromContext(ctx, " cattegories_api")
	time.Sleep(0 * time.Second)
	defer span.Finish()
	return generateRandomArray(5)
}

func generateRandomArray(num int) []string {
	randomArray := make([]string, num)
	for i := 0; i < num; i++ {
		randomArray[i] = fmt.Sprintf("%d", i)
	}
	return randomArray
}
