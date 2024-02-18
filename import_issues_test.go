package main

import (
	"encoding/csv"
	"os"
	"strings"
	"testing"
)

func TestGetEnvToken(t *testing.T) {
	expectedToken := "test_token"
	err := setEnvironmentVariable("GITHUB_TOKEN", expectedToken)
	if err != nil {
		t.Fatal("Failed to set environment variable:", err)
	}
	defer os.Unsetenv("GITHUB_TOKEN")

	token := getEnvToken("GITHUB_TOKEN")
	if token != expectedToken {
		t.Fatalf("getEnvToken = %q; want %q", token, expectedToken)
	}
}
func TestReadRecord(t *testing.T) {
	csvContent := `Title1,Description1
Title2,Description2
,,`
	r := csv.NewReader(strings.NewReader(csvContent))

	// Test first record
	record, err := readRecord(r)
	if err != nil {
		t.Fatal("readRecord failed:", err)
	}
	if len(record) != 2 || record[0] != "Title1" || record[1] != "Description1" {
		t.Fatalf("readRecord = %q; want %q", record, []string{"Title1", "Description1"})
	}

	// Test second record
	record, err = readRecord(r)
	if err != nil {
		t.Fatal("readRecord failed:", err)
	}
	if len(record) != 2 || record[0] != "Title2" || record[1] != "Description2" {
		t.Fatalf("readRecord = %q; want %q", record, []string{"Title2", "Description2"})
	}

	// Test invalid record
	_, err = readRecord(r)
	if err == nil {
		t.Fatal("readRecord should have failed with invalid record")
	}
}

// setEnvironmentVariable is a helper function to set an environment variable for testing.
func setEnvironmentVariable(key, value string) error {
	return os.Setenv(key, value)
}
