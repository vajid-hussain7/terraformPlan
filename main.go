package main

import (
	"encoding/json"
	"fmt"
	"jsonTest/model"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func main() {

	CreateCommit()

	// gitFile, err := http.Get("https://raw.githubusercontent.com/vajid-hussain7/terraformPlan/main/terraform-plan.json")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// gitFileData, err := io.ReadAll(gitFile.Body)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// var json interface{}
	// var MainField map[string]interface{}

	// // extract data from json

	// // file, err := ExtractFile()
	// // if err != nil {
	// // 	log.Fatalln("file reading lead to error ", err)
	// // }
	// // create and write the terraform plan data to plan.json file
	// err = os.WriteFile("plan.json", gitFileData, os.FileMode(0744))
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println("created")

	// // unmarshel file
	// err = Unmarshel(gitFileData, &json)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// // create the json result for storing in to db
	// MainField = json.(map[string]interface{})
	// result := FinalResult(MainField["resource_changes"].([]interface{}))

	// err = db.InitDB()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// db.InsertResourse(result)

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

func CreateCommit() {

	localRepo, err := git.PlainOpen("./")
	if err != nil {
		fmt.Println(err)
	}

	// _, err = localRepo.Branch("test")
	// if err != nil {
	// 	err := localRepo.CreateBranch(&config.Branch{Name: "test"})
	// 	if err != nil {
	// 		fmt.Println("create a branch ", err)
	// 	}
	// } else {
	// 	fmt.Println(err)
	// }

	workTree, err := localRepo.Worktree()
	if err != nil {
		fmt.Println("from worktree", err)
	}

	branch := plumbing.NewBranchReferenceName("test1")

	err = workTree.Checkout(&git.CheckoutOptions{Branch: branch, Create: true})
	if err != nil {
		fmt.Println(err)
	}

	_, err = workTree.Add("./")
	if err != nil {
		fmt.Println(err)
	}

	_, err = workTree.Commit("test1 branch first commit ", &git.CommitOptions{})
	if err != nil {
		fmt.Println(err)
	}

	var auth = &http.BasicAuth{
		Username: "vajid-hussain7",
		Password: os.Getenv("git_token"),
	}

	err = localRepo.Push(&git.PushOptions{RemoteName: "origin", Auth: auth, Force: true})
	if err != nil {
		fmt.Println(err)
	}

}

// ErrBranchExists
