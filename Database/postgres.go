package database

import (
	"database/sql"
	"example/baseProject/envvariable"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
)

var (
	host     = envvariable.Host
	port     = envvariable.Port
	user     = envvariable.User
	password = envvariable.Password
	dbname   = envvariable.DBName
)

// const (
// 	host     =
// 	port     = 5432
// 	user     = "postgres"
// 	password = "1234"
// 	dbname   = "test"
// )

type Postgres struct {
	db *sql.DB
}

func NewPostgres() *Postgres {
	//TODO:SingleTone Pattern

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	dbconn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println("Error couldn't connect to Database:", err)
		return nil
	}

	return &Postgres{
		db: dbconn,
	}
}

func (obj *Postgres) Insert(tableName string, argsKeys []string, argsVals []string) error {
	sql := fmt.Sprintf(`INSERT INTO %s(%s) VALUES ('%s')`, tableName, strings.Join(argsKeys, ","), strings.Join(argsVals, `','`))
	fmt.Println(sql)
	_, err := obj.db.Query(sql)
	if err != nil {
		log.Println("Error inserting into Database:", err)
		return err
	}
	return nil
}

func (obj *Postgres) SelectById(tableName string, id int) (*sql.Rows, error) {
	sql := fmt.Sprintf("SELECT id, name,price FROM %s  WHERE id = %d", tableName, id)
	fmt.Println(sql)
	row, err := obj.db.Query(sql)
	if err != nil {
		log.Printf("Error SELECT from Database table %s ERROR Massage %s", tableName, err)
		return nil, err
	}
	return row, nil
}

func (obj *Postgres) CloseDB() error {
	err := obj.db.Close()
	if err != nil {
		log.Println("Error While Clossing Database:", err)
		return err
	}
	fmt.Println("Closed the database connection")
	return nil
}
