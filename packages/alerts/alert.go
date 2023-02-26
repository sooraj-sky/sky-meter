package skyalerts

import (
	"bytes"
	"encoding/json"
	gomail "gopkg.in/gomail.v2"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	models "sky-meter/models"

	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"strconv"
)

func SendMail(i models.SmtpErr) {

	log.Println(i.Mailto, "mailto")

	emailPass := os.Getenv("emailpass")
	if emailPass == "" {
		log.Fatal("Please specify the emailpass as environment variable, e.g. env emailpass=your-pass go run http-server.go")
	}

	t := template.New("error.html")

	var err error
	t, err = t.ParseFiles("templates/error.html")
	if err != nil {
		log.Println(err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, i); err != nil {
		log.Println(err)
	}

	log.Println("emails areee", i.Mailto)
	for k := range i.Mailto {

		result := tpl.String()
		m := gomail.NewMessage()
		m.SetHeader("From", os.Getenv("emailFrom"))
		m.SetHeader("To", i.Mailto[k])
		m.SetHeader("Subject", i.Subject)
		m.SetBody("text/html", result)
		intPort, _ := strconv.Atoi(os.Getenv("EmailPort"))

		d := gomail.NewDialer(os.Getenv("emailServer"), intPort, os.Getenv("emailFrom"), emailPass)

		if err := d.DialAndSend(m); err != nil {
			log.Println(err)

		}

	}
}

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
	opsgenieSecret := os.Getenv("opsgeniesecret")
	opsgenieSecretString := "GenieKey " + opsgenieSecret
	if opsgenieSecret == "" {
		log.Fatal("Please specify the opsgeniesecret as environment variable, e.g. export opsgeniesecret=")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", opsgenieSecretString)
	resp, err := apiclient.Do(req)
	if err != nil {
		log.Println(err.Error())
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
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
