package main

import (
	"encoding/json"
	"fmt"
	"io"
	"jsonTest/db"
	"jsonTest/model"
	"log"
	"net/http"
	"os"
)

func main() {

	gitFile, err := http.Get("https://raw.githubusercontent.com/vajid-hussain7/terraformPlan/main/terraform-plan.json")
	if err != nil {
		log.Fatalln(err)
	}

	gitFileData, err := io.ReadAll(gitFile.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var json interface{}
	var MainField map[string]interface{}

	// extract data from json

	// file, err := ExtractFile()
	// if err != nil {
	// 	log.Fatalln("file reading lead to error ", err)
	// }
	// create and write the terraform plan data to plan.json file
	err = os.WriteFile("plan.json", gitFileData, os.FileMode(0744))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("created")

	// unmarshel file
	err = Unmarshel(gitFileData, &json)
	if err != nil {
		log.Fatalln(err)
	}

	// create the json result for storing in to db
	MainField = json.(map[string]interface{})
	result := FinalResult(MainField["resource_changes"].([]interface{}))

	err = db.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	db.InsertResourse(result)

}

// unmarshel the json to empty interface
func Unmarshel(data []byte, model *interface{}) error {
	return json.Unmarshal(data, model)
}

// read file from project parent root
func ExtractFile() ([]byte, error) {
	return os.ReadFile("terraform-plan.json")
}

// create the final result
func FinalResult(resousrseChange []interface{}) model.ResourseChanges {

	var finalResult = model.ResourseChanges{}

	for _, val := range resousrseChange {

		// conversions && creating the end result
		eachResoure := val.(map[string]interface{})
		change := eachResoure["change"].(map[string]interface{})
		action := change["actions"].([]interface{})

		finalResult.Resourses = append(finalResult.Resourses, model.Resourses{Action: action[0].(string), Type: eachResoure["type"].(string)})
	}

	return finalResult
}


func CreateCommit(){
	
}