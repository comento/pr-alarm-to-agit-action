package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
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

func getApiGithubEvent(event string, token string, private bool) GitPullRequestUrl {
	client := &http.Client{}
	req, err := http.NewRequest("GET", event, nil)
	if err != nil {
		panic(err)
	}

	if private {
		req.Header.Add("Authorization", "Bearer "+token)
	}
	req.Header.Add("Content-Type", "application/json")

	res, _ := client.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var gitPullRequestUrl GitPullRequestUrl
	json.Unmarshal([]byte(body), &gitPullRequestUrl)

	return gitPullRequestUrl
}

func getApiGithubCommitsEvent(url string, token string, private bool) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", "Bearer " + token)
	req.Header.Add("Content-Type", "application/json")

	res, _ := client.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var gitCommitsApi []GitCommitsApi
	json.Unmarshal([]byte(body), &gitCommitsApi)

	var commitHistory string
	for i := range gitCommitsApi {
		commitHistory += "*" + gitCommitsApi[i].GitCommit.Committer.Name + "*: " + strings.Replace(gitCommitsApi[i].GitCommit.Message, "\n", "\\n", -1) + "\\n"
	}

	return commitHistory
}

func main() {
	url := os.Getenv("INPUT_URL")
	event := os.Getenv("INPUT_EVENT")
	private, _ := strconv.ParseBool(os.Getenv("INPUT_PRIVATE"))
	token := os.Getenv("INPUT_TOKEN")

	var gitPullRequestUrl GitPullRequestUrl
	gitPullRequestUrl = getApiGithubEvent(event, token, private)

	var commitHistory = getApiGithubCommitsEvent(gitPullRequestUrl.CommitsUrl, token, private)

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
