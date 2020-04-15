package main

import (
	"fmt"
	_ "database/sql"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type env struct {
	host, user, password, port, name string
}

func init() {
	if err := godotenv.Load(); err != nil {
		panic("No .env file found")
	}
}

func main() {
	fmt.Println(getEnv())
}

func getEnv() env {
	dbEnv := env{
		host: checkEnv("host"),
		user: checkEnv("user"),
		password: checkEnv("password"),
		port: checkEnv("port"),
		name: checkEnv("name"),
	}
	return dbEnv
}

func checkEnv(str string) string {
	str, exists := os.LookupEnv("DB_" + strings.ToUpper(str))
	if !exists {
		panic(".env not properly configured")
	}
	return str
}
