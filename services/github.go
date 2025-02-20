package services

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

type Commit struct {
	Sha    string `json:"sha"`
	Commit struct {
		Message string `json:"message"`
	} `json:"commit"`
}

func FetchLatestCommit(owner, repo string) (*Commit, error) {
	client := resty.New()
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits", owner, repo)
	resp, err := client.R().Get(url)

	if err != nil {
		log.Println("❌ GitHub API request failed:", err)
		return nil, err
	}

	var commits []Commit
	err = json.Unmarshal(resp.Body(), &commits)
	if err != nil || len(commits) == 0 {
		log.Println("❌ Failed to parse commits:", err)
		return nil, err
	}

	return &commits[0], nil
}
