package main

import (
	"fmt"
	"database/sql"
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
	dbEnv := getEnv()
	psqlConnectionStr := fmt.Sprintf(
		"host=%s port=%s user=%s password='%s' dbname=%s sslmode=disable",
		dbEnv.host, dbEnv.port, dbEnv.user, dbEnv.password, dbEnv.name,
	)
	db, err := sql.Open("postgres", psqlConnectionStr)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to postgreSQL!")
	db.Close()
}

func getEnv() env {
	return env{
		host: checkEnv("host"),
		user: checkEnv("user"),
		password: checkEnv("password"),
		port: checkEnv("port"),
		name: checkEnv("name"),
	}
}

func checkEnv(str string) string {
	str, exists := os.LookupEnv("DB_" + strings.ToUpper(str))
	if !exists {
		panic(".env is missing environment variable: '" + str + "'")
	}
	return str
}
