package main

import (
	"DonTaskMe-backend/internal/db"
	"DonTaskMe-backend/routing"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	_ = godotenv.Load()
	db.InitDb()
	defer db.Disconnect()

	var mode string
	if len(os.Args) > 1 && os.Args[1] == "--prod" {
		mode = "release"
	} else {
		mode = "debug"
	}

	server := routing.GetServer(mode)
	err := server.Run()
	if err != nil {
		log.Fatalln("Couldn't start the server: ", err)
	}
}
