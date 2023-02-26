package yamlops

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path/filepath"
	models "sky-meter/models"
)

// This is a Go function named InputYml that returns a models.UserInput object.
func InputYml() models.UserInput {

	// Get the absolute path of the settings.yml file in the current directory.
	filename, _ := filepath.Abs("./settings.yml")

	// Read the contents of the settings.yml file into a byte slice.
	yamlFile, err := ioutil.ReadFile(filename)

	// If an error occurs while reading the file, log it.
	if err != nil {
		log.Println(err)
	}

	// Create a new UserInput object to unmarshal the YAML into.
	var config models.UserInput

	// Unmarshal the YAML data into the UserInput object.
	err = yaml.Unmarshal(yamlFile, &config)

	// If an error occurs while unmarshaling, log it.
	if err != nil {
		log.Println(err)
	}

	// Return the UserInput object.
	return config
}
