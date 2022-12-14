package htttpserver

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	api "github.com/sooraj-sky/sky-meter/packages/api"
)

func InitServer() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Please specify the HTTP port as environment variable, e.g. env PORT=8081 go run http-server.go")
	}

	log.Println("listening on port", port)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", api.HomeLink)
	router.HandleFunc("/health", api.SelfStatusLink)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
