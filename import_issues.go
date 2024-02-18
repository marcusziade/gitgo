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

const initialBackoff = 5 * time.Second

func main() {
	ctx := context.Background()

	token := getEnvToken("GITHUB_TOKEN")
	if token == "" {
		return
	}

	client := createGitHubClient(ctx, token)
	processCSVFile(ctx, client, "tasks.csv")
}

func getEnvToken(envVar string) string {
	token := os.Getenv(envVar)
	if token == "" {
		fmt.Printf("The %s environment variable is not set.\n", envVar)
	}
	return token
}

func createGitHubClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func processCSVFile(ctx context.Context, client *github.Client, filePath string) {
	csvfile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening the CSV file:", err)
		return
	}
	defer csvfile.Close()

	r := csv.NewReader(csvfile)
	backoff := initialBackoff

	for {
		record, err := readRecord(r)
		if err != nil {
			break
		}
		createIssueFromRecord(ctx, client, record, &backoff)
	}
	fmt.Println("All issues created successfully.")
}

func readRecord(r *csv.Reader) ([]string, error) {
	record, err := r.Read()
	if err == io.EOF {
		return nil, err
	}
	if err != nil {
		fmt.Println("Error reading the CSV file:", err)
		return nil, err
	}
	if len(record) < 2 || strings.TrimSpace(record[0]) == "" {
		return nil, fmt.Errorf("invalid record")
	}
	return record, nil
}

func createIssueFromRecord(ctx context.Context, client *github.Client, record []string, backoff *time.Duration) {
	issueRequest := &github.IssueRequest{
		Title: github.String(strings.TrimSpace(record[0])),
		Body:  github.String(strings.TrimSpace(record[1])),
	}

	for {
		created, shouldRetry := createGitHubIssue(ctx, client, issueRequest, *backoff)
		if created || !shouldRetry {
			break
		}
		*backoff *= 2
	}
	*backoff = initialBackoff
}

func createGitHubIssue(ctx context.Context, client *github.Client, issueRequest *github.IssueRequest, backoff time.Duration) (bool, bool) {
	fmt.Printf("Creating issue: %s", *issueRequest.Title)

	issue, resp, err := client.Issues.Create(ctx, "marcusziade", "vibify", issueRequest)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusForbidden {
			fmt.Printf("\rHit rate limit, backing off for %d seconds\n", backoff)
			time.Sleep(backoff * time.Second)
			return false, true
		}
		fmt.Printf("\rError creating issue for '%s': %s\n", *issueRequest.Title, err)
		return false, false
	}

	fmt.Printf("\rSuccessfully created issue '%s': %s\n", *issueRequest.Title, *issue.HTMLURL)
	time.Sleep(2 * time.Second)
	fmt.Println()
	return true, false
}
