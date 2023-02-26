package htttpserver

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	api "sky-meter/packages/api"
	skyenv "sky-meter/packages/env"
)

func InitServer() {
	allEnv := skyenv.GetEnv()
	port := allEnv.Port
	log.Println("listening on port", port)
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", api.HomeLink)
	router.HandleFunc("/health", api.SelfStatusLink)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
