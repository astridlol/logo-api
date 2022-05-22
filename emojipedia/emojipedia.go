package emojipedia

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	ErrNoEmoji = errors.New("no emoji found")
	ErrNoUrl   = errors.New("no url found")
)

func Search(term string, emojiType string) ([]byte, error) {
	searchRes, err := http.Get(fmt.Sprintf("https://emojipedia.org/search/?q=%s", term))

	if err != nil {
		return nil, err
	}

	if searchRes.StatusCode < 200 || searchRes.StatusCode > 299 {
		return nil, errors.New(fmt.Sprintf("search request failed with status code: %d", searchRes.StatusCode))
	}

	defer searchRes.Body.Close()

	searchDoc, err := goquery.NewDocumentFromReader(searchRes.Body)

	if err != nil {
		return nil, err
	}

	pageUrl, exists := searchDoc.Find(`ol.search-results li h2 a`).Attr("href")

	if !exists {
		return nil, ErrNoEmoji
	}

	emojiRes, err := http.Get(fmt.Sprintf("https://emojipedia.org%s", pageUrl))

	if err != nil {
		return nil, err
	}

	if emojiRes.StatusCode < 200 || emojiRes.StatusCode > 299 {
		return nil, errors.New(fmt.Sprintf("emoji page request failed with status code: %d", emojiRes.StatusCode))
	}

	defer emojiRes.Body.Close()

	emojiDoc, err := goquery.NewDocumentFromReader(emojiRes.Body)

	if err != nil {
		return nil, err
	}

	emojiUrl, exists := emojiDoc.Find(`section.vendor-list ul li div.vendor-container div.vendor-image img`).Attr("src")

	switch emojiType {
	// No need to do an Apple statement, as it's like that by default
	case "android":
		{
			emojiUrl = strings.Replace(emojiUrl, "apple/325/", "google/313/", 3)
		}
	case "discord":
		{
			emojiUrl = strings.Replace(emojiUrl, "apple/325/", "twitter/322/", 3)
		}
	}

	if !exists {
		return nil, ErrNoUrl
	}

	imageRes, err := http.Get(emojiUrl)

	if err != nil {
		return nil, err
	}

	if imageRes.StatusCode < 200 || imageRes.StatusCode > 299 {
		return nil, errors.New(fmt.Sprintf("image download request failed with status code: %d", imageRes.StatusCode))
	}

	body, _ := ioutil.ReadAll(imageRes.Body)

	return body, nil
}
