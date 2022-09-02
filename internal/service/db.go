package service

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	Name                    = "myFirstDatabase"
	UsersCollectionName     = "users"
	WorkspaceCollectionName = "workspaces"
	ListCollectionName      = "lists"
)

var (
	Client *mongo.Client
	DB     *mongo.Database
)

func InitDb() {
	log.Println("Initialising Mongo Database")
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("ATLAS_URI")))

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	log.Printf("Connecting to: %s", Name)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database connection established!")
	Client = client
	DB = client.Database(Name)
}

func Disconnect() {
	_ = Client.Disconnect(context.TODO())
}
