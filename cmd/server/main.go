package main

import (
	"context"
	"log"
	"net/http"

	"graphQlDemo"
	"graphQlDemo/ent"
	"graphQlDemo/ent/migrate"
	"graphQlDemo/auth"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/lib/pq"
	"os"
	"github.com/joho/godotenv"
)

func main() {
	godotenvErr := godotenv.Load()
	if godotenvErr != nil {
	  log.Fatal("Error loading .env file")
	}
	connStr := os.Getenv("DATABASE")
	client, clientErr := ent.Open("postgres", connStr)
	if clientErr != nil {
		log.Fatal("opening ent client", clientErr)
	}
	defer client.Close()

	if err := client.Schema.Create(
		context.Background(),
		migrate.WithGlobalUniqueID(true),
	); err != nil {
		log.Fatal("running schema migration", err)
	}

	srv := handler.NewDefaultServer(graphQlDemo.NewSchema(client))
	router := http.NewServeMux()
	router.Handle("/", playground.Handler("GraphQL Playground", "/graphql"))
	router.Handle("/graphql", auth.Middleware(client, srv))
	server := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	log.Println("server started on http://localhost:8081")
	log.Fatal(server.ListenAndServe())
}
