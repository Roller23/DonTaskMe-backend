package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

const (
	Name                = "myFirstDatabase"
	UsersCollectionName = "users"
)

var (
	Client *mongo.Client
)

func InitDb() {
	log.Println("Initialising Mongo Database")
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("ATLAS_URI")))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	log.Printf("Connecting to: %s", os.Getenv("DB_NAME"))
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database connection established!")
	Client = client
}

func Disconnect() {
	_ = Client.Disconnect(context.TODO())
}
