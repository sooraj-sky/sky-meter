package jsonops

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	models "github.com/sooraj-sky/sky-meter/models"
)

func InputJson() models.JsonInput {
	jsonFile, err := os.Open("input.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var endpoints models.JsonInput
	json.Unmarshal(byteValue, &endpoints)
	return endpoints
}
