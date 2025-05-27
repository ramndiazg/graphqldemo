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
	if err := client.Schema.Create(
		context.Background(),
		migrate.WithGlobalUniqueID(true),
	); err != nil {
		log.Fatal("running schema migration", err)
	}

	srv := handler.NewDefaultServer(graphQlDemo.NewSchema(client))

	authHandler := auth.Middleware(client, srv)

	http.Handle("/",
		playground.Handler("GraphQL Demo", "/query"),
	)
	http.Handle("/query", authHandler)

	log.Println("listening on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("http server terminated", err)
	}
}
