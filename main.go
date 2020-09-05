package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/shomali11/slacker"
	"github.com/joho/godotenv"
)

func main() {
	botKey := getVarFromENV("SLACK_BOT_TOKEN")
	bot := slacker.NewClient(botKey)

	definition := &slacker.CommandDefinition{
		Description: "Get job postings from greenhouse boards",
		Example: 	 "tarantula crawl github",
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			board := request.Param("word")
			postings, err := getJobPostings("https://boards.greenhouse.io/" + board)
			if err != nil {
				log.Fatal(err)
			}
			response.Reply(postings, slacker.WithThreadReply(true))
		},
	}

	bot.Command("tarantula crawl <word>", definition)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
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
