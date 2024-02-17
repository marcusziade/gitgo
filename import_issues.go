package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

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
	for {
		record, err := r.Read()
		if err != nil {
			break
		}

		issueRequest := &github.IssueRequest{
			Title: github.String(strings.TrimSpace(record[0])),
			Body:  github.String(strings.TrimSpace(record[1])),
		}

		issue, _, err := client.Issues.Create(ctx, "marcusziade", "vibify", issueRequest)
		if err != nil {
			fmt.Printf("Error creating issue for '%s': %s\n", *issueRequest.Title, err)
			continue
		}

		fmt.Printf("Successfully created issue '%s': %s\n", *issueRequest.Title, *issue.HTMLURL)
	}
}
