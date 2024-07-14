package main

import (
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func HttpToByte(url string) []byte {
	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to fetch URL: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Fatalf("Failed to fetch URL, status code: %d", response.StatusCode)
	}

	imageBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	return imageBytes
}

func GetProfilePicture(username string) ([]byte, error) {
	url := "https://www.instagram.com/" + username
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Fatalf("Failed to fetch page, status code: %d", response.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	var profilePictureURL string
	doc.Find("meta[property='og:image']").Each(func(i int, s *goquery.Selection) {
		content, exists := s.Attr("content")
		if exists {
			profilePictureURL = content
		}
	})

	if profilePictureURL != "" {
		// fmt.Println("Profile Picture URL:", profilePictureURL)
		return HttpToByte(profilePictureURL), nil
	} else {
		// fmt.Println("Profile Picture URL not found")
		return nil, errors.New("profile Picture URL not found")
	}

}
