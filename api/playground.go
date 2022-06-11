package handler

import (
	"github.com/99designs/gqlgen/graphql/playground"
	"log"
	"net/http"
	"os"
)

var playgroundUrl = os.Getenv("PLAYGROUND_ENDPOINT")
var playgroundServer = playground.Handler("GraphQL playground", playgroundUrl)

func Playground(w http.ResponseWriter, r *http.Request) {
	log.Printf("New Playground request from %s", r.UserAgent())
	for name, values := range r.Header {
		// Loop over all values for the name.
		for _, value := range values {
			log.Printf("Header: %s - Name: %s\n", name, value)
		}
	}
	playgroundServer.ServeHTTP(w, r)
}
