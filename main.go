package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"jsonTest/model"
	"log"
	internalHttp "net/http"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/v63/github"
	"golang.org/x/oauth2"
)

func main() {

	CreateCommit()

	CreatePullRequest()

	// download file from github server contet
	gitFile, err := internalHttp.Get("https://raw.githubusercontent.com/vajid-hussain7/terraformPlan/main/terraform-plan.json")
	if err != nil {
		log.Fatalln(err)
	}

	// read from the response body
	gitFileData, err := io.ReadAll(gitFile.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// create and write the terraform plan data to plan.json file
	err = os.WriteFile("plan.json", gitFileData, os.FileMode(0744))
	if err != nil {
		fmt.Println(err)
		return
	}

	// var json interface{}
	// var MainField map[string]interface{}

	// // extract data from json file
	// // file, err := ExtractFile()
	// // if err != nil {
	// // 	log.Fatalln("file reading lead to error ", err)
	// // }

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

	// find the local repository from project root eg:=".git"
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

	// fetch the worktree from local repo
	workTree, err := localRepo.Worktree()
	if err != nil {
		fmt.Println("from worktree", err)
	}

	// define branch
	branch := plumbing.NewBranchReferenceName("test1")

	// create branch on workTree, if branch exist return error
	err = workTree.Checkout(&git.CheckoutOptions{Branch: branch, Create: true})
	if err != nil {
		fmt.Println(err)
	}

	// add changes to state state include all changes "."
	_, err = workTree.Add(".")
	if err != nil {
		fmt.Println(err)
	}

	// create a new commit with stages changes along with commit message
	_, err = workTree.Commit("test pr branch second commit ", &git.CommitOptions{})
	if err != nil {
		fmt.Println(err)
	}

	// create auth, it help on git push
	var auth = &http.BasicAuth{
		Username: "vajid-hussain7",
		Password: os.Getenv("git_token"),
	}

	// push the commit to remote origin
	err = localRepo.Push(&git.PushOptions{RemoteName: "origin", Auth: auth, Force: true})
	if err != nil {
		fmt.Println(err)
	}

}

// create a pull request to main branch
func CreatePullRequest() {
	// initalize a parent content
	ctx := context.Background()

	// prerequirement auth
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("git_token")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	owner := "vajid-hussain7"
	repo := "terraformPlan"
	headBranch := "test1"
	baseBranch := "main"

	// create a pr instance
	pr := &github.NewPullRequest{
		Title: github.String("test1 pr2 title"),
		Head:  github.String(headBranch),
		Base:  github.String(baseBranch),
		Body:  github.String("Description of test1 pr2"),
	}

	// create a pull request 
	prResp, _, err := client.PullRequests.Create(ctx, owner, repo, pr)
	if err != nil {
		log.Fatalf("Error creating pull request: %v", err)
	}

	fmt.Printf("Pull request created: %s\n", *prResp.HTMLURL)
}
