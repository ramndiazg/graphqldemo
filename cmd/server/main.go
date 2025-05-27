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
)

const (
	connStr = "host=localhost port=5432 user=admin dbname=tools-back password=1234 sslmode=disable"
)

func main() {
	client, err := ent.Open("postgres", connStr)
	if err != nil {
		log.Fatal("opening ent client", err)
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
