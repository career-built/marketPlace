package utilities

import (
	"fmt"
	"os"
	"strconv"
)

const (
	hostKey     = "DB_HOST"
	portKey     = "DB_PORT"
	userKey     = "DB_USER"
	passwordKey = "DB_PASSWORD"
	dbnameKey   = "DB_NAME"
)

var (
	Host     = getEnvVar(hostKey, "localhost")
	Port     = getEnvVarInt(portKey, 5432)
	User     = getEnvVar(userKey, "postgres")
	Password = getEnvVar(passwordKey, "1234")
	DBName   = getEnvVar(dbnameKey, "test")
)

func getEnvVar(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func getEnvVarInt(key string, defaultValue int) int {
	strValue := getEnvVar(key, "")
	if strValue == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(strValue)
	if err != nil {
		fmt.Printf("Invalid integer value for %s, using default %d\n", key, defaultValue)
		return defaultValue
	}

	return value
}
