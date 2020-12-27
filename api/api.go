// Package api configures the api based on environment (dev, test, prod)
package api

import (
	"fmt"
	"log"
	"os"

	"github.com/denisecase/go-hunt-sql/api/router"
	"github.com/denisecase/go-hunt-sql/api/seeder"
	"github.com/joho/godotenv"
)

var appServer = router.Server{}

// init function is private - it starts with a lowercase letter -
// init() reads values from .env into the system
func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("Cannot read from .env file in root folder.")
	}
}

// Run function is public - it starts with an upper case letter -
// Run() preloads data and starts listening
func Run() {
	fmt.Println("Starting server")

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	} else {
		fmt.Println("Reading env values")
	}
	appPort := os.Getenv("PORT")
	appEnv := os.Getenv("ENV")
	devDbDriver := os.Getenv("DEV_DB_DRIVER")
	devDB := os.Getenv("DEV_DB")
	devDBInMemory := os.Getenv("DEV_DB_INMEMORY")

	dbDriver := os.Getenv("DB_DRIVER")
	dbUsername := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	appServer.Initialize(dbDriver, dbUsername, dbPassword, dbPort, dbHost, dbName, appEnv, devDbDriver, devDB, devDBInMemory)

	seeder.Load(appServer.DB)

	appServer.Run(":" + appPort)

}
