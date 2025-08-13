package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ecetinerdem/forseer/api"
	"github.com/ecetinerdem/forseer/database"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	db, err := database.NewDB()

	if err != nil {
		log.Fatal("Database connection error: ", err)
	}

	defer db.Close()

	database.RunMigrations(db)

	server := api.NewServer(db)
	PORT := os.Getenv("PORT")
	log.Println("Server starting on the designated port")
	log.Fatal(http.ListenAndServe(":"+PORT, server.Router))

}
