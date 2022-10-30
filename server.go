package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jasonlvhit/gocron"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	dbops "sky-meter/dbops"
	httpreponser "sky-meter/httpres"
)

type AllEndpoints []struct {
	URL       string `json:"url",omitempty`
	Timeout   int    `json:"timeout",omitempty`
	SkipSsl   bool   `json:"skip_ssl",omitempty`
	Frequency uint64 `json:"frequency",omitempty`
	Group     string `json:"group",omitempty`
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to sky-meter")
}

func getStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	httpresdata, _ := httpreponser.GetHttpdata("https://bing.com")
	w.Write(httpresdata)
	return
}

func httpSyntheticCheck(endpoint string, time uint64) {
	gocron.Every(time).Second().Do(callEndpoint, endpoint)
	<-gocron.Start()
	fmt.Println(time)
}

func callEndpoint(endpoint string) {
	httpresdata, _ := httpreponser.GetHttpdata(endpoint)
	fmt.Println(string(httpresdata))
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

	for i := 0; i < len(endpoints); i++ {
		dbops.InsertSearchUrl(endpoints[i].URL, endpoints[i].Timeout, endpoints[i].SkipSsl, endpoints[i].Frequency, endpoints[i].Group)
		httpSyntheticCheck(endpoints[i].URL, endpoints[i].Frequency)
	}

	fmt.Println("listening on port 8080")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/stats", getStats).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))

}