package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Println("The GITHUB_TOKEN environment variable is not set.")
		return
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	csvfile, err := os.Open("tasks.csv")
	if err != nil {
		fmt.Println("Error opening the CSV file:", err)
		return
	}
	defer csvfile.Close()

	r := csv.NewReader(csvfile)
	var backoff time.Duration = 5

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading the CSV file:", err)
			break
		}
		if len(record) < 2 || strings.TrimSpace(record[0]) == "" {
			continue
		}

		issueRequest := &github.IssueRequest{
			Title: github.String(strings.TrimSpace(record[0])),
			Body:  github.String(strings.TrimSpace(record[1])),
		}

		fmt.Printf("Creating issue: %s", strings.TrimSpace(record[0]))

		issue, resp, err := client.Issues.Create(ctx, "marcusziade", "vibify", issueRequest)
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusForbidden {
				fmt.Printf("\rHit rate limit, backing off for %d seconds\n", backoff)
				time.Sleep(backoff * time.Second)
				backoff *= 2
			} else {
				fmt.Printf("\rError creating issue for '%s': %s\n", *issueRequest.Title, err)
			}
			continue
		}
		backoff = 1
		fmt.Printf("\rSuccessfully created issue '%s': %s\n", *issueRequest.Title, *issue.HTMLURL)

		// Visual loading state for 2 seconds
		time.Sleep(2 * time.Second)
		fmt.Println()
	}
	fmt.Println("All issues created successfully.")
}
