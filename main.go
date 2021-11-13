package main

import (
	"DonTaskMe-backend/graph"
	"DonTaskMe-backend/graph/generated"
	"DonTaskMe-backend/internal/db"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	_ = godotenv.Load()
	port := os.Getenv("PORT")
	db.InitDb()
	defer db.Disconnect()
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	log.Printf("GraphiQL http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
