package main

import (
	"fmt"
	"log"
	"movie-review/api/middleware"
	"movie-review/api/repository"
	"movie-review/config"
	"movie-review/db"
	"movie-review/graph"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
)

func main() {
	err := config.LoadConfig(".config")
	if err != nil {
		panic(fmt.Sprintf("cannot load config: %v", err))
	}
	db := db.Connect()
	defer db.Close()
	repos := repository.InitRepositories(db)
	router := chi.NewRouter()
	router.Use(middleware.HandleCORS)
	router.Use(middleware.AddRepoToContext(repos))
	
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.NewRootResolvers(db)))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", config.ConfigVal.Port)
	log.Fatal(http.ListenAndServe(":"+config.ConfigVal.Port, router))
}
