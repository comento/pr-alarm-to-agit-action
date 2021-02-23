package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type GitPullRequestUrl struct {
	Title string `json:"title"`
	HtmlUrl string `json:"html_url"`
	CommitsUrl string `json:"commits_url"`
}

type GitCommitsApi struct {
	GitCommit GitCommit `json:"commit"`
}

type GitCommit struct {
	Committer Committer `json:"committer"`
	Message   string    `json:"message"`
}

type Committer struct {
	Name string `json:"name"`
}

func getApiGithubEvent(event string) GitPullRequestUrl {
	res, err := http.Get(event)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var gitPullRequestUrl GitPullRequestUrl
	json.Unmarshal([]byte(body), &gitPullRequestUrl)

	return gitPullRequestUrl
}

func getApiGithubCommitsEvent(url string) string {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var gitCommitsApi []GitCommitsApi
	json.Unmarshal([]byte(body), &gitCommitsApi)

	var commitHistory string
	for i := range gitCommitsApi {
		commitHistory += gitCommitsApi[i].GitCommit.Committer.Name + ": " + gitCommitsApi[i].GitCommit.Message + "\\n"
	}

	return commitHistory
}

func main() {
	url := os.Getenv("INPUT_URL")
	event := os.Getenv("INPUT_EVENT")

	var gitPullRequestUrl GitPullRequestUrl
	gitPullRequestUrl = getApiGithubEvent(event)

	var commitHistory = getApiGithubCommitsEvent(gitPullRequestUrl.CommitsUrl)

	payload := strings.NewReader("{\"text\": \"PR OPEN: " + gitPullRequestUrl.Title + "\\n" + gitPullRequestUrl.HtmlUrl + "\\n\\n" + commitHistory + "\"}")

	// http post
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "Application/json")

	res, _ := client.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
}
