package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const GITHUB_API = "https://api.github.com"

type PullRequest struct {
	Number   int    `json:"number"`
	Title    string `json:"title"`
	MergedAt string `json:"merged_at"`
	HtmlURL  string `json:"html_url"`
	RepoName string `json:"repo_name"`
}

func checkMergedStatus(owner, repo string, prNumber int) (string, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls/%d", GITHUB_API, owner, repo, prNumber)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var prData struct {
		MergedAt *string `json:"merged_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&prData); err != nil {
		return "", err
	}

	if prData.MergedAt == nil {
		return "", nil 
	}
	return *prData.MergedAt, nil
}


func getUserPRs(username string) ([]PullRequest, error) {
	url := fmt.Sprintf("%s/search/issues?q=author:%s+type:pr", GITHUB_API, username)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Items []struct {
			Number  int    `json:"number"`
			Title   string `json:"title"`
			HtmlURL string `json:"html_url"`
			RepoURL string `json:"repository_url"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var prs []PullRequest
	for _, pr := range result.Items {
		
		repoParts := strings.Split(pr.RepoURL, "/")
		owner := repoParts[len(repoParts)-2] 
		repoName := repoParts[len(repoParts)-1]

		
		mergedAt, err := checkMergedStatus(owner, repoName, pr.Number)
		if err != nil {
			fmt.Println("Error checking merge status:", err)
		}

		prs = append(prs, PullRequest{
			Number:   pr.Number,
			Title:    pr.Title,
			MergedAt: mergedAt,
			HtmlURL:  pr.HtmlURL,
			RepoName: repoName,
		})
	}

	return prs, nil
}

func main() {
	r := gin.Default()

	r.GET("/prs/:username", func(c *gin.Context) {
		username := c.Param("username")

		prs, err := getUserPRs(username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching PRs"})
			return
		}

		c.JSON(http.StatusOK, prs)
	})

	log.Fatal(r.Run(":3000"))
}
