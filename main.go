package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func main() {

	ScrapeSingleImage("www.example.com")

	fileName := "test.jpg"
	URL := ""
	err := downloadImage(URL, fileName)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("File %s downloaded in current working directory", fileName)
}

// ScrapeSingleImage -  Image Download
func ScrapeSingleImage(s string) string {
	res, err := http.Get(s)
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

	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if property, _ := s.Attr("property"); property == "og:image" {
			content, _ := s.Attr("content")
			fmt.Printf(content)

			content = URL
		}
	})
	return URL
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

	return nil
}
