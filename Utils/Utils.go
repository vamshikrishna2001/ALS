package Utils

import (
	"encoding/json"
	"fmt"
	"os"
)

func CreateFile(Filename string, object interface{}) error {
	file, err := os.Create(Filename)
	if err != nil {
		fmt.Println("Error creating JSON file:", err)
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(object, " ", "  ") // Indented for better readability
	if err != nil {
		fmt.Println("Error marshaling struct to JSON:", err)
	}

	// Write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
	return err

}
