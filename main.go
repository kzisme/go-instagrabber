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

type InstagramUserData struct {
	Graphql struct {
		User struct {
			Biography              string      `json:"biography"`
			BlockedByViewer        bool        `json:"blocked_by_viewer"`
			BusinessCategoryName   interface{} `json:"business_category_name"`
			CategoryEnum           interface{} `json:"category_enum"`
			CategoryName           interface{} `json:"category_name"`
			ConnectedFbPage        interface{} `json:"connected_fb_page"`
			CountryBlock           bool        `json:"country_block"`
			EdgeFelixVideoTimeline struct {
				Count    int64         `json:"count"`
				Edges    []interface{} `json:"edges"`
				PageInfo struct {
					EndCursor   interface{} `json:"end_cursor"`
					HasNextPage bool        `json:"has_next_page"`
				} `json:"page_info"`
			} `json:"edge_felix_video_timeline"`
			EdgeFollow struct {
				Count int64 `json:"count"`
			} `json:"edge_follow"`
			EdgeFollowedBy struct {
				Count int64 `json:"count"`
			} `json:"edge_followed_by"`
			EdgeMediaCollections struct {
				Count    int64         `json:"count"`
				Edges    []interface{} `json:"edges"`
				PageInfo struct {
					EndCursor   interface{} `json:"end_cursor"`
					HasNextPage bool        `json:"has_next_page"`
				} `json:"page_info"`
			} `json:"edge_media_collections"`
			EdgeMutualFollowedBy struct {
				Count int64         `json:"count"`
				Edges []interface{} `json:"edges"`
			} `json:"edge_mutual_followed_by"`
			EdgeOwnerToTimelineMedia struct {
				Count int64 `json:"count"`
				Edges []struct {
					Node struct {
						Typename             string `json:"__typename"`
						AccessibilityCaption string `json:"accessibility_caption"`
						CommentsDisabled     bool   `json:"comments_disabled"`
						Dimensions           struct {
							Height int64 `json:"height"`
							Width  int64 `json:"width"`
						} `json:"dimensions"`
						DisplayURL  string `json:"display_url"`
						EdgeLikedBy struct {
							Count int64 `json:"count"`
						} `json:"edge_liked_by"`
						EdgeMediaPreviewLike struct {
							Count int64 `json:"count"`
						} `json:"edge_media_preview_like"`
						EdgeMediaToCaption struct {
							Edges []struct {
								Node struct {
									Text string `json:"text"`
								} `json:"node"`
							} `json:"edges"`
						} `json:"edge_media_to_caption"`
						EdgeMediaToComment struct {
							Count int64 `json:"count"`
						} `json:"edge_media_to_comment"`
						EdgeMediaToTaggedUser struct {
							Edges []interface{} `json:"edges"`
						} `json:"edge_media_to_tagged_user"`
						FactCheckInformation   interface{} `json:"fact_check_information"`
						FactCheckOverallRating interface{} `json:"fact_check_overall_rating"`
						GatingInfo             interface{} `json:"gating_info"`
						ID                     string      `json:"id"`
						IsVideo                bool        `json:"is_video"`
						Location               struct {
							HasPublicPage bool   `json:"has_public_page"`
							ID            string `json:"id"`
							Name          string `json:"name"`
							Slug          string `json:"slug"`
						} `json:"location"`
						MediaOverlayInfo interface{} `json:"media_overlay_info"`
						MediaPreview     string      `json:"media_preview"`
						Owner            struct {
							ID       string `json:"id"`
							Username string `json:"username"`
						} `json:"owner"`
						SharingFrictionInfo struct {
							BloksAppURL               interface{} `json:"bloks_app_url"`
							ShouldHaveSharingFriction bool        `json:"should_have_sharing_friction"`
						} `json:"sharing_friction_info"`
						Shortcode          string `json:"shortcode"`
						TakenAtTimestamp   int64  `json:"taken_at_timestamp"`
						ThumbnailResources []struct {
							ConfigHeight int64  `json:"config_height"`
							ConfigWidth  int64  `json:"config_width"`
							Src          string `json:"src"`
						} `json:"thumbnail_resources"`
						ThumbnailSrc string `json:"thumbnail_src"`
					} `json:"node"`
				} `json:"edges"`
				PageInfo struct {
					EndCursor   string `json:"end_cursor"`
					HasNextPage bool   `json:"has_next_page"`
				} `json:"page_info"`
			} `json:"edge_owner_to_timeline_media"`
			EdgeRelatedProfiles struct {
				Edges []interface{} `json:"edges"`
			} `json:"edge_related_profiles"`
			EdgeSavedMedia struct {
				Count    int64         `json:"count"`
				Edges    []interface{} `json:"edges"`
				PageInfo struct {
					EndCursor   interface{} `json:"end_cursor"`
					HasNextPage bool        `json:"has_next_page"`
				} `json:"page_info"`
			} `json:"edge_saved_media"`
			ExternalURL            interface{} `json:"external_url"`
			ExternalURLLinkshimmed interface{} `json:"external_url_linkshimmed"`
			Fbid                   string      `json:"fbid"`
			FollowedByViewer       bool        `json:"followed_by_viewer"`
			FollowsViewer          bool        `json:"follows_viewer"`
			FullName               string      `json:"full_name"`
			HasArEffects           bool        `json:"has_ar_effects"`
			HasBlockedViewer       bool        `json:"has_blocked_viewer"`
			HasChannel             bool        `json:"has_channel"`
			HasClips               bool        `json:"has_clips"`
			HasGuides              bool        `json:"has_guides"`
			HasRequestedViewer     bool        `json:"has_requested_viewer"`
			HighlightReelCount     int64       `json:"highlight_reel_count"`
			ID                     string      `json:"id"`
			IsBusinessAccount      bool        `json:"is_business_account"`
			IsJoinedRecently       bool        `json:"is_joined_recently"`
			IsPrivate              bool        `json:"is_private"`
			IsVerified             bool        `json:"is_verified"`
			OverallCategoryName    interface{} `json:"overall_category_name"`
			ProfilePicURL          string      `json:"profile_pic_url"`
			ProfilePicURLHd        string      `json:"profile_pic_url_hd"`
			RequestedByViewer      bool        `json:"requested_by_viewer"`
			RestrictedByViewer     interface{} `json:"restricted_by_viewer"`
			ShouldShowCategory     bool        `json:"should_show_category"`
			Username               string      `json:"username"`
		} `json:"user"`
	} `json:"graphql"`
	LoggingPageID           string      `json:"logging_page_id"`
	ProfilePicEditSyncProps interface{} `json:"profile_pic_edit_sync_props"`
	ShowFollowDialog        bool        `json:"show_follow_dialog"`
	ShowSuggestedProfiles   bool        `json:"show_suggested_profiles"`
	ShowViewShop            bool        `json:"show_view_shop"`
	ToastContentOnLoad      interface{} `json:"toast_content_on_load"`
}

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
