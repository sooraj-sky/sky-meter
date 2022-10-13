package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	httpreponser "sky-meter/httpres"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to sky-meter")
}

func getStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	httpresdata, _ := httpreponser.GetHttpdata("https://apple.com")
	w.Write(httpresdata)
	return
}

func main() {
	fmt.Println("listening on port 8080")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/stats", getStats).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
