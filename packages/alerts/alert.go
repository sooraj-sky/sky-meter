package skyalerts

import (
	"bytes"
	"encoding/json"
	gomail "gopkg.in/gomail.v2"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	models "sky-meter/models"
	skyenv "sky-meter/packages/env"

	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
	"strconv"
)

// SendMail function sends an email to the specified recipient(s) with the given email content

func SendMail(i models.SmtpErr) {

	// Get the environment variables for the email server details and the password for the email account.
	allEnv := skyenv.GetEnv()
	emailPass := allEnv.EmailPass

	// Create a new template from the error.html file and handle any errors that occur during the parsing of the template.
	t := template.New("error.html")
	var err error
	t, err = t.ParseFiles("templates/error.html")
	if err != nil {
		log.Println(err)
	}

	// Execute the template with the SmtpErr instance to get the final HTML email body.
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, i); err != nil {
		log.Println(err)
	}

	// Loop through each recipient email address listed in the SmtpErr instance, create a new email message, set the necessary headers,
	// and send the email using the SMTP server details and the email account password.
	for k := range i.Mailto {

		result := tpl.String()
		m := gomail.NewMessage()
		m.SetHeader("From", allEnv.EmailFrom)
		m.SetHeader("To", i.Mailto[k])
		m.SetHeader("Subject", i.Subject)
		m.SetBody("text/html", result)

		// Convert email port string to integer
		intPort, _ := strconv.Atoi(allEnv.EmailPort)

		// Create a dialer for the email server and authenticate using email credentials
		d := gomail.NewDialer(allEnv.EmailServer, intPort, allEnv.EmailFrom, emailPass)

		// Dial the email server and send the email
		if err := d.DialAndSend(m); err != nil {
			log.Println(err)

		}

	}
}

type error interface {
	Error() string
}

// This function takes in the URL of the endpoint that caused the error, an error object with the error details, and a group name as strings, and creates a new alert in OpsGenie.
func OpsgenieCreateAlert(errorurl string, description error, group string) string {

	// Create a message for the alert and get the environment variable for the OpsGenie API key.
	downMessege := "Alert Endpint " + errorurl + " is Down"
	allEnv := skyenv.GetEnv()
	opsgenieSecret := allEnv.OpsgenieSecret

	// Create an OpsGenie alert client with the specified API key
	alertClient, err := alert.NewClient(&client.Config{
		ApiKey: opsgenieSecret,
	})

	// Create a new alert with the specified details
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

	// Return the request ID of the created alert
	return createResult.RequestId

}

// CheckAlertStatus retrieves the status of an OpsGenie alert with the specified request ID
func CheckAlertStatus(alertRequestId string) string {

	// Create a new HTTP client.
	apiclient := &http.Client{}

	// Build the URL for the API call using the alert request ID.
	url := "https://api.opsgenie.com/v2/alerts/requests/" + alertRequestId + "?identifierType=id"

	// Retrieve the Opsgenie secret key from the environment variables.
	allEnv := skyenv.GetEnv()
	opsgenieSecret := allEnv.OpsgenieSecret
	opsgenieSecretString := "GenieKey " + opsgenieSecret

	// Create a new GET request with the Opsgenie secret key.
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", opsgenieSecretString)

	// Send the HTTP request to Opsgenie API.
	resp, err := apiclient.Do(req)
	if err != nil {
		log.Println(err.Error())
	}
	defer resp.Body.Close()

	// Read the response body and unmarshal it into a struct.
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
	}
	var responseObject models.OpsGenieAlertStatus
	json.Unmarshal(bodyBytes, &responseObject)

	// Create a new alert client using the Opsgenie secret key.
	alertClient, err := alert.NewClient(&client.Config{
		ApiKey: opsgenieSecret,
	})

	// Retrieve the alert status using the alert ID from the response.
	getResult, err := alertClient.Get(nil, &alert.GetAlertRequest{
		IdentifierType:  alert.ALERTID,
		IdentifierValue: responseObject.Data.AlertID,
	})

	if err != nil {
		log.Printf("error: %s\n", err)
	}

	// Return the alert status.
	return getResult.Status
}
