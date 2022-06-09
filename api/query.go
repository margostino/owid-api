package handler

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/margostino/owid-api/fetcher"
	"github.com/margostino/owid-api/graph"
	"github.com/margostino/owid-api/graph/generated"
	"net/http"
)

var server = newServer()

func Query(w http.ResponseWriter, r *http.Request) {
	server.ServeHTTP(w, r)
	//http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	//http.Handle("/query", server)
	//log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	//log.Fatal(http.ListenAndServe(":"+port, nil))
}

func newServer() *handler.Server {
	c := generated.Config{Resolvers: &graph.Resolver{Fetcher: fetcher.Fetch}}
	return handler.NewDefaultServer(generated.NewExecutableSchema(c))
}
