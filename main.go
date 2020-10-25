package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	err := downloadImage(ScrapeSingleImage("https://www.instagram.com/p/Bh2mQ1UBcJS/"))

	if err != nil {
		log.Fatal(err)
	}
}

// ScrapeSingleImage -  Image Download
func ScrapeSingleImage(singleImageURL string) (string, string) {
	res, err := http.Get(singleImageURL)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var URL string
	var fileName string

	// Returning CDN URL from generic Instagram post
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if property, _ := s.Attr("property"); property == "og:image" {
			content, _ := s.Attr("content")

			URL = content
		}
	})

	// Returning Post-Identifier for filename
	doc.Find("meta").Each(func(i int, sel *goquery.Selection) {
		if property, _ := sel.Attr("property"); property == "og:url" {
			myURL, err := url.Parse(singleImageURL)
			if err != nil {
				log.Fatal(err)
			}

			fileName = path.Base(myURL.Path) + ".jpg"
		}
	})

	fmt.Printf("Filename: %s CDN URL: %s\n", fileName, URL)
	return URL, fileName
}

func downloadImage(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}
	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	fmt.Printf("File %s downloaded in current working directory", fileName)

	return nil
}
