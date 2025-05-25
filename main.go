package main

import (
	"fmt"
	"log"

	"workship-disaster-api/db"
	"workship-disaster-api/router"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// เชื่อมต่อ PostgreSQL
	dbConn, err := db.ConnectPostgres()
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	// Run database migrations
	if err := db.RunMigrations(dbConn); err != nil {
		log.Fatal("Error running migrations:", err)
	}

	// เชื่อมต่อ Redis
	rdb, err := db.ConnectRedis()
	if err != nil {
		log.Fatal(err)
	}
	defer rdb.Close()

	// สร้าง API
	r := router.SetupRouter(dbConn, rdb)

	fmt.Println("Server is running on port 8080")
	r.Run(":8080")
}
