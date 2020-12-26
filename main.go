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

func main() {
	// Check every 5 minutes
	for range time.NewTicker(5 * time.Minute).C {
		webScraper()
	}
}

func webScraper() {
	// Instantiate default collector
	c := colly.NewCollector()

	var htmlElementClean string

	// On every a element which has href attribute call callback
	c.OnHTML("div[id='availability']", func(e *colly.HTMLElement) {
		htmlElement := e.Text
		htmlElementClean = strings.Replace(htmlElement, "\n", "", -1)
		// fmt.Printf("Availability: %q", htmlElementClean)
	})
	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		t := time.Now().UTC()
		fmt.Println(t.String(), "visiting:", r.URL.String())
	})

	// Start scraping
	c.Visit("https://www.amazon.co.uk/dp/B08H97NYGP/ref=twister_B08J4RCVXW?_encoding=UTF8&psc=1")

	// Email me if PLaystation is available
	if htmlElementClean[:22] != "Currently unavailable." {
		fmt.Println("PS5 Available.")
		mail(htmlElementClean)
	}
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
