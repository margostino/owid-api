package handler

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/margostino/owid-api/fetcher"
	"github.com/margostino/owid-api/graph"
	"github.com/margostino/owid-api/graph/generated"
	"github.com/margostino/owid-api/metrics"
	"net/http"
)

var server = newServer()

func Query(w http.ResponseWriter, r *http.Request) {
	go metrics.PublishRequest(r)
	server.ServeHTTP(w, r)
}

func newServer() *handler.Server {
	c := generated.Config{Resolvers: &graph.Resolver{Fetcher: fetcher.Fetch}}
	return handler.NewDefaultServer(generated.NewExecutableSchema(c))
}
