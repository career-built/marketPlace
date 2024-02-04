package product

import (
	"database/sql"
	"log"
	"strconv"
)

type DbConnector interface {
	// start() error
	Insert(tableName string, argsKeys []string, argsVals []string) error
	CloseDB() error

	//TODO: Remove SQL Package Make it puplic
	SelectById(tableName string, id int) (*sql.Rows, error)
}

type MessageBroker interface {
	PublishMessages(exchange string, queueName string, messages []string) error
	ConsumeMessages(queueName string, handler func(message string)) error
	Close() error
}

type ProductService struct {
	dbConnector   DbConnector
	messageBroker MessageBroker
}

func NewProductService(connector DbConnector, message_Broker MessageBroker) *ProductService {
	return &ProductService{
		dbConnector:   connector,
		messageBroker: message_Broker,
	}
}

func (obj *ProductService) Add(product *Product, queueName string) error {
	priceStr := strconv.Itoa(product.PRICE)
	keys := []string{"name", "price"}
	vals := []string{
		product.NAME,
		priceStr,
	}
	err := obj.dbConnector.Insert("product", keys, vals)
	if err != nil {
		log.Println("Product not inserted", err)
		return err
	}
	// Notify the Other Service that there an update in the product list
	message := "Product:: " + product.NAME + " With Price:: " + priceStr + " inserted"
	queueName = "productlist"
	err = obj.messageBroker.PublishMessages("", queueName, []string{message})

	if err != nil {
		log.Println("PublishMessages Not Sent", err)
		return err
	}

	return nil
}

func (obj *ProductService) GetByID(id int) (*Product, error) {
	row, err := obj.dbConnector.SelectById("product", id)
	if err != nil {
		log.Printf("Error while retrive the product with %d , Erorr:,%s ", id, err)
		return nil, err
	}
	product := &Product{}

	// Check if there's a row to scan
	if row.Next() {
		// Scan the values from the row into the Product struct
		err := row.Scan(&product.PRICE, &product.NAME)
		if err == nil {
			log.Printf("Error scanning row values: %s", err)
			return nil, err
		}
		return product, nil
	} else {
		// No rows found
		return nil, nil
	}

}
