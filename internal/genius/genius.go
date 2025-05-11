package genius

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type Genius struct{}

func NewGenius() *Genius {
	return &Genius{}
}

func (g *Genius) GetTrackPageURL(title string) (string, error) {
	url := "https://genius.com/search?q=" + title

	browser := rod.New().Timeout(10 * time.Second)
	if err := browser.Connect(); err != nil {
		return "", errors.New("failed to connect to browser: " + err.Error())
	}
	defer browser.Close()

	page, err := browser.Page(proto.TargetCreateTarget{URL: url})
	if err != nil {
		return "", errors.New("failed to create page: " + err.Error())
	}

	if err := page.WaitLoad(); err != nil {
		return "", errors.New("failed to load page: " + err.Error())
	}

	html, err := page.HTML()
	if err != nil {
		return "", errors.New("failed to get HTML: " + err.Error())
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Printf("Failed to parse HTML: %v", err)
		log.Println("Opening search page", url)
		return url, nil
	}

	href, exists := doc.Find("a.mini_card").First().Attr("href")
	if !exists {
		log.Println("No exact match found, opening search page:", url)
		return url, nil
	}

	log.Println("Found exact match:", href)
	return href, nil
}
