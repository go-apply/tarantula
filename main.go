package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	postings, err := GetJobPostings("https://boards.greenhouse.io/github")
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Website:")
	fmt.Printf(postings)
}

// GetJobPostings gets the latest jobs given and returns them as a list
func GetJobPostings(url string) (string, error) {

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
