package yamlops

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	models "sky-meter/models"
)

func InputYml() models.UserInput {
	filename, _ := filepath.Abs("./settings.yml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Println(err)
	}

	var config models.UserInput

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Println(err)
	}
	os.Setenv("emailFrom", config.Email[0].Sender)
	os.Setenv("emailServer", config.Email[0].Server)
	os.Setenv("EmailPort", string(config.Email[0].Port))

	return config
}
