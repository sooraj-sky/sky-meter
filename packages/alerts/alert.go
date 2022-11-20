package skyalerts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sky-meter/models"

	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
)

type error interface {
	Error() string
}

func OpsgenieCreateAlert(errorurl string, description error, group string) string {
	downMessege := "Alert Endpint " + errorurl + " is Down"
	opsgenieSecret := os.Getenv("opsgeniesecret")
	if opsgenieSecret == "" {
		log.Fatal("Please specify the opsgeniesecret as environment variable, e.g. export opsgeniesecret=")
	}

	alertClient, err := alert.NewClient(&client.Config{
		ApiKey: opsgenieSecret,
	})

	createResult, _ := alertClient.Create(nil, &alert.CreateAlertRequest{
		Message:     downMessege,
		Description: description.Error(),
		Tags:        []string{"P1", errorurl},
		Details: map[string]string{
			"Group": group,
		},
		Priority: alert.P1,
	})

	if err != nil {
		log.Printf("error: %s\n", err)
	}

	return createResult.RequestId

}

func CheckAlertStatus(alertRequestId string) string {
	apiclient := &http.Client{}
	url := "https://api.opsgenie.com/v2/alerts/requests/" + alertRequestId + "?identifierType=id"
	log.Println(url)
	opsgenieSecret := os.Getenv("opsgeniesecret")
	opsgenieSecretString := "GenieKey " + opsgenieSecret
	if opsgenieSecret == "" {
		log.Fatal("Please specify the opsgeniesecret as environment variable, e.g. export opsgeniesecret=")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Print(err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", opsgenieSecretString)
	resp, err := apiclient.Do(req)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	var responseObject models.OpsGenieAlertStatus
	json.Unmarshal(bodyBytes, &responseObject)

	alertClient, err := alert.NewClient(&client.Config{
		ApiKey: opsgenieSecret,
	})

	getResult, err := alertClient.Get(nil, &alert.GetAlertRequest{
		IdentifierType:  alert.ALERTID,
		IdentifierValue: responseObject.Data.AlertID,
	})

	if err != nil {
		log.Printf("error: %s\n", err)
	}

	return getResult.Status
}
