package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"nathanielwheeler.com/controllers"
	"nathanielwheeler.com/models"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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
	// Start up database connection
	dbEnv := getEnv()
	psqlConnectionStr := fmt.Sprintf(
		"host=%s port=%s user=%s password='%s' dbname=%s sslmode=disable",
		dbEnv.host, dbEnv.port, dbEnv.user, dbEnv.password, dbEnv.name,
	)

	// Initialize services
	ss, err := models.NewSubsService(psqlConnectionStr)
	if err != nil {
		panic(err)
	}
	defer ss.Close()
	ss.AutoMigrate()

	// Initialize controllers
	staticC := controllers.NewStatic()
	subsC := controllers.NewSubscribers()

	// Route Handling
	router := mux.NewRouter()
	router.Handle("/", staticC.Home).Methods("GET")
	router.Handle("/resume", staticC.Resume).Methods("GET")
	router.HandleFunc("/subscribe", subsC.New).Methods("GET")
	router.HandleFunc("/subscribe", subsC.Create).Methods("POST")

	// Start that server!
	http.ListenAndServe(":3000", router)
}


// #region DB HELPERS

func getEnv() env {
	return env{
		host:     checkEnv("host"),
		user:     checkEnv("user"),
		password: checkEnv("password"),
		port:     checkEnv("port"),
		name:     checkEnv("name"),
	}
}

func checkEnv(str string) string {
	str, exists := os.LookupEnv("DB_" + strings.ToUpper(str))
	if !exists {
		panic(".env is missing environment variable: '" + str + "'")
	}
	return str
}

// #endregion