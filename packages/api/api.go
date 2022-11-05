package api

import (
	"fmt"
	"net/http"
)

func HomeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to sky-meter")
}

func SelfStatusLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status OK")
}
