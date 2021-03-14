package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/PuerkitoBio/goquery"
)

func main() {

	//https://www.instagram.com/{username}/?__a=1

	// Initial Call Of First 12
	// https://www.instagram.com/graphql/query/?query_id=17888483320059182&&variables={"id":"375193502","first":12}

	// Subsequent Calls
	// https://www.instagram.com/graphql/query/?query_id=17888483320059182&&variables={"id":"375193502","first":12,"after":"END_CURSOR_STRING_HERE"}

	// The function can be simplified further by using the following:
	// https://www.instagram.com/graphql/query/?query_id=17888483320059182&&variables={"id":"375193502","first":12}

	PaginateAndDownload(GetUserJson("kzita93"))

	doc, err := goquery.NewDocument("https://www.instagram.com/kzita93/")
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("a[href]").Each(func(index int, item *goquery.Selection) {
		href, _ := item.Attr("href")
		fmt.Printf("link: %s - anchor text: %s\n", href, item.Text())
	})

	//////err := downloadImage(ScrapeSingleImage(""))

	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func PaginateAndDownload(userJson string) {
	userData := InstagramUserData{}
	jsonData := json.Unmarshal([]byte(userJson), &userData)

	if jsonData != nil {
		log.Fatal(jsonData)
	}

	// This is the number of user posts - initial call holds 0-11 + end_cursor
	fmt.Println(userData.Graphql.User.EdgeOwnerToTimelineMedia.Count)
}

func GetUserJson(userName string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.instagram.com/"+userName+"/channel/?__a=1", nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	//log.Println(string(body))

	return string(body)
}

// ScrapeSingleImage -  Image Download
func ScrapeSingleImage(singleImageURL string) (string, string, string) {
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
	var authorID string

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

	// Pull authorID from page to create directory
	doc.Find("meta").Each(func(i int, sel *goquery.Selection) {
		if property, _ := sel.Attr("property"); property == "instapp:owner_user_id" {
			content, _ := sel.Attr("content")

			authorID = content
		}
	})

	fmt.Printf("Filename: %s\n CDN URL: %s\n", fileName, URL)
	return URL, fileName, authorID
}

func downloadImage(URL, fileName string, authorID string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}

	newpath := filepath.Join(filepath.Dir(os.Args[0]), authorID)
	os.Mkdir(newpath, os.ModePerm)

	//Create a empty file
	file, err := os.Create(filepath.Join(newpath, fileName))
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
