package handler

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/margostino/owid-api/graph"
	"github.com/margostino/owid-api/graph/generated"
	"net/http"
)

var server = handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

func Query(w http.ResponseWriter, r *http.Request) {
	server.ServeHTTP(w, r)
	//http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	//http.Handle("/query", server)
	//log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	//log.Fatal(http.ListenAndServe(":"+port, nil))
}
