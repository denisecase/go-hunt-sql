package router

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/denisecase/go-hunt-sql/api/models"
)

// Server struct contains database and router
type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

// Initialize sets up the datastore based on environment variables
func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName, AppEnv, DevDbDriver, DevDB, DevDBInMemory string) {

	var err error

	fmt.Println("AppEnv=", AppEnv)
	fmt.Println("DevDbDriver=", DevDbDriver)
	fmt.Println("DevDBInMemory=", DevDBInMemory)

	if AppEnv == "development" && DevDbDriver == "sqlite3" && DevDBInMemory == "true" {
		dbString := "development sqlite3 in-memory database"
		fmt.Println("Initializing ", dbString)
		server.DB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
		if err != nil {
			fmt.Printf("Cannot connect to %s\n", dbString)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("Connected to the %s\n", dbString)
		}
		server.DB.Exec("PRAGMA foreign_keys = ON")
	} else if AppEnv == "development" && DevDbDriver == "sqlite3" && DevDBInMemory != "true" {
		dbString := "development sqlite3 database"
		fmt.Println("Initializing ", dbString)
		server.DB, err = gorm.Open(sqlite.Open("gorm-hunt.db"), &gorm.Config{})
		if err != nil {
			fmt.Printf("Cannot connect to %s\n", dbString)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("Connected to the %s\n", dbString)
		}
		server.DB.Exec("PRAGMA foreign_keys = ON")
	} else {
		dbString := "postgres database"
		fmt.Println("Initializing ", dbString)
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=CST", DbHost, DbUser, DbPassword, DbName, DbPort)
		server.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			fmt.Printf("Cannot connect to %s", dbString)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("Connected to the %s", dbString)
		}
	}

	server.DB.AutoMigrate(&models.User{}, &models.Team{}) //database migration

	server.Router = mux.NewRouter()

	server.initializeRoutes()
}

// Run starts the server / begins listening
func (server *Server) Run(port string) {
	fmt.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(port, server.Router))
}
