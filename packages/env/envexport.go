package env

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	models "sky-meter/models"
	yamlops "sky-meter/packages/yamlops"
)

// To add a new env variable, add that variable to models.AllEnvs struct.

// This function is used to get the keys of a struct type. It does this by creating a pointer to an instance of the models.AllEnvs struct type, using reflection to get the type of this instance, and then checking that the type is a struct.
// If it is a struct, it iterates through all of the fields of the struct and appends their names to a slice of strings, which it returns at the end.
func GetEnvStructKeys() (envKeys []string) {
	a := &models.AllEnvs{}
	t := reflect.TypeOf(*a)
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			envKeys = append(envKeys, t.Field(i).Name)

		}
	}
	return envKeys
}

// This function is used to initialize the environment variables needed by the program.
// It does this by calling the GetEnvStructKeys function to get a slice of the keys of the models.AllEnvs struct, and then iterating through this slice.
// For each key, it retrieves the corresponding environment variable using the os.Getenv function, and checks that it is not an empty string.
// If the environment variable is empty, it logs an error message indicating that the variable needs to be specified and terminates the program.

func InitEnv() {

	envNames := GetEnvStructKeys()

	var envStatus []string

	endpoints := yamlops.InputYml()

	if endpoints.Opegenie.Enabled != true {
		os.Setenv("OpsgenieSecret", "tmp-opsgenie-key")
	}
	if endpoints.Email.Enabled != true {
		os.Setenv("EmailPass", "tmp-key")
		os.Setenv("EmailFrom", "tmp-key")
		os.Setenv("EmailPort", "tmp-key")
		os.Setenv("EmailServer", "tmp-key")
	}

	for i := range envNames {
		envValue := os.Getenv(envNames[i])
		if envValue == "" {
			envStatus = append(envStatus, "\n Variable "+envNames[i]+" is not set \n")
		}
	}
	if envStatus != nil {

		log.Fatal(envStatus)

	}

}

// This function is used to retrieve and return the environment variables needed by the program as an instance of the models.AllEnvs struct.
// It does this by calling the GetEnvStructKeys function to get a slice of the keys of the models.AllEnvs struct, creating an empty map[string]string and then iterating through the slice of keys.
// For each key, it retrieves the corresponding environment variable using the os.Getenv function and stores it in the map using the key as the map key.
// Next, it marshals the map into a JSON string using the json.Marshal function.
// It then unmarshals the JSON string into an instance of the models.AllEnvs struct using the json.Unmarshal function, and returns this instance.
// If either the json.Marshal or json.Unmarshal functions return an error, it logs the error.

func GetEnv() (outputEnv models.AllEnvs) {
	envNames := GetEnvStructKeys()

	var Envvalues models.AllEnvs
	mapOfEnv := make(map[string]string)

	for i := range envNames {
		mapOfEnv[envNames[i]] = os.Getenv(envNames[i])
	}

	jsonStr, err := json.Marshal(mapOfEnv)

	if err != nil {
		log.Println(err)
	}

	if err := json.Unmarshal(jsonStr, &Envvalues); err != nil {
		log.Println(err)
	}

	return Envvalues
}
