package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	dbops "sky-meter/dbops"
	httpreponser "sky-meter/httpres"
)

type AllEndpoints struct {
	Endpoints []struct {
		URL     string `json:"url"`
		Timeout int    `json:"timeout"`
		SkipSsl bool   `json:"skip_ssl"`
	} `json:"endpoints"`
}

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

	jsonFile, err := os.Open("input.json")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened input.json")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var endpoints AllEndpoints

	json.Unmarshal(byteValue, &endpoints)

	for i := 0; i < len(endpoints.Endpoints); i++ {
		dbops.InsertSearchUrl(endpoints.Endpoints[i].URL,endpoints.Endpoints[i].Timeout,endpoints.Endpoints[i].SkipSsl)
	}

	fmt.Println("listening on port 8080")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/stats", getStats).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
