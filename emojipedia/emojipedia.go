package emojipedia

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

var (
	ErrNoEmoji = errors.New("no emoji found")
	ErrNoUrl   = errors.New("no url found")
)

func Search(term string) ([]byte, error) {
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
