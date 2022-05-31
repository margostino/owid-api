package handler

import (
	"github.com/99designs/gqlgen/graphql/playground"
	"net/http"
)

var playgroundServer = playground.Handler("GraphQL playground", "/query")

func Playground(w http.ResponseWriter, r *http.Request) {
	playgroundServer.ServeHTTP(w, r)
}
