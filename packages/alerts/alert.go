package skyalerts

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type error interface {
	Error() string
}

func OpsgenieCreateAlert(errorurl string, description error) {
	downMessege := "Alert Endpint " + errorurl + " is Down"
	url := "https://api.opsgenie.com/v2/alerts"
	values := map[string]string{"message": downMessege, "description": description.Error(), "priority": "P1"}
	json_data, err := json.Marshal(values)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_data))
	opsgenieSecret := os.Getenv("opsgeniesecret")
	if opsgenieSecret == "" {
		log.Fatal("Please specify the opsgeniesecret as environment variable, e.g. export opsgeniesecret=")
	}
	ApiSec := "GenieKey " + opsgenieSecret
	req.Header.Set("Authorization", ApiSec)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	log.Println("response Status:", resp.Status)
}
