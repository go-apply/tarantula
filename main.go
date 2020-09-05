package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/slack-go/slack"
	"github.com/joho/godotenv"
)

func main() {
	slackApiKey := getVarFromENV("slack")

	postings, err := getJobPostings("https://boards.greenhouse.io/github")
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Website:")
	fmt.Printf(postings)
}

func getVarFromENV(key string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

// getJobPostings gets the latest jobs given and returns them as a list
func getJobPostings(url string) (string, error) {

	// Get the HTML
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	// Convert HTML into goquery document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	// Save each .opening as a list
	openings := ""
	doc.Find(".opening").Each(func(i int, s *goquery.Selection) {
		temp := strings.Trim(s.Text(), "\n \t")
		openings += "- " + temp
	})
	return openings, nil
}
