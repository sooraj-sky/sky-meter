package jsonops

import (
	"os"
	models "sky-meter/models"
	"fmt"
	"io/ioutil"
	"encoding/json"

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
