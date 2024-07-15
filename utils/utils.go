package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func LoadJSON(filename string, v interface{}) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func SaveJSON(filename string, v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
