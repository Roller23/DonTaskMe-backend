package main

import (
	"DonTaskMe-backend/routing"
	"context"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Couldn't load .env:", err)
	}

	mongoURI := os.Getenv("ATLAS_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancelDB := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelDB()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalln("Could not disconnect from mongoDB: ", err)
	}

	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Fatalln("Couldn't disconnect properly: ", err)
		}
	}(client, ctx)

	var mode string
	if len(os.Args) > 1 && os.Args[1] == "--prod" {
		mode = "release"
	} else {
		mode = "debug"
	}

	server := routing.GetServer(mode)
	err = server.Run()
	if err != nil {
		log.Fatalln("Couldn't start the server: ", err)
	}
}
