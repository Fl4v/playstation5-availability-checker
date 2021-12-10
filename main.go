package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

var siteUrl string

func main() {
	// Check on start
	webScraper("div[id='availability']", "https://www.amazon.co.uk/dp/B08H97NYGP/ref=twister_B08J4RCVXW?_encoding=UTF8&psc=1", "buy")
	// Then check every n minute(s)
	for range time.NewTicker(2 * time.Minute).C {
		availalbe := webScraper("div[id='availability']", "https://www.amazon.co.uk/dp/B08H97NYGP/ref=twister_B08J4RCVXW?_encoding=UTF8&psc=1", "unavailable")
		fmt.Println("sss")
		// When Playstation is available, end loop
		if availalbe {
			mail(siteUrl)
		}
	}
}

func webScraper(div string, site string, availability string) bool {
	// Instantiate default collector
	c := colly.NewCollector()

	var htmlElementClean string

	// On every div element with id availability
	c.OnHTML(div, func(e *colly.HTMLElement) {
		htmlElement := e.Text
		htmlElementClean = strings.Replace(htmlElement, "\n", "", -1)
	})

	c.OnRequest(func(r *colly.Request) {
		t := time.Now().UTC()
		fmt.Println(t.String(), "visiting:", r.URL.String())
	})

	// Start scraping
	c.Visit(site)

	// Receive email if Playstation is available
	if !strings.Contains(htmlElementClean, availability) {
		fmt.Println("PS5 Available.")
		siteUrl = site
		return true
	}
	return false
}

func mail(body string) {
	// Set up authentication information.
	auth := smtp.PlainAuth("", os.Getenv("SMTP_EMAIL"), os.Getenv("SMTP_PASSWORD"), "smtp.gmail.com")

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	to := []string{os.Getenv("SMTP_EMAIL")}
	msg := []byte("To: " + os.Getenv("SMTP_EMAIL") + "\r\n" +
		"Subject: Playstation 5 Alert!\r\n" +
		"\r\n" +
		body +
		"\r\n")

	err := smtp.SendMail("smtp.gmail.com:25", auth, os.Getenv("SMTP_EMAIL"), to, msg)
	if err != nil {
		log.Fatal(err)
	}
}
