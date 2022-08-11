package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	PhotoAPI = "https://api.pexels.com/v1"
	VideoAPI = "https://api.pexels.com/videos/"
)

type Client struct {
	Token          string
	hc             http.Client
	RemainingTimes int32
}

func NewClient(token string) *Client {
	c := http.Client{}
	return &Client{Token: token, hc: c}
}

type SearchResult struct {
	Page         int32   `json:"page"`
	PerPage      int32   `json:"per_page"`
	TotalResults int32   `json:"total_results"`
	NextPage     string  `json:"next_page"`
	Photos       []Photo `json:"photos"`
}

type CuratedResult struct {
	Page     int32   `json:"page"`
	PerPage  int32   `json:"per_page"`
	NextPage string  `json:"next_page"`
	Photos   []Photo `json:"photos"`
}

type Photo struct {
	Id              int32       `json:"id"`
	Width           int32       `json:"width"`
	Height          int32       `json:"height"`
	Url             string      `json:"url"`
	Photographer    string      `json:"photographer"`
	PhotographerUrl string      `json:"photographer_url"`
	Src             PhotoSource `json:"src"`
}

type PhotoSource struct {
	Original  string `json:"original"`
	Large     string `json:"large"`
	Large2x   string `json:"large2x"`
	Medium    string `json:"medium"`
	Small     string `json:"small"`
	Portrait  string `json:"portrait"`
	Square    string `json:"square"`
	Landscape string `json:"landscape"`
	Tiny      string `json:"tiny"`
}

type VideoSearchResult struct {
	Page         int32   `json:"page"`
	PerPage      int32   `json:"per_page"`
	TotalResults int     `json:"total_results"`
	NextPage     string  `json:"next_page"`
	Videos       []Video `json:"videos"`
}
type Video struct {
	Id            int32           `json:"id"`
	Width         int32           `json:"width"`
	Height        int32           `json:"height"`
	Url           string          `json:"url"`
	Image         string          `json:"image"`
	FullRes       interface{}     `json:"full_res"`
	Duration      float64         `json:"duration"`
	VideoFiles    []VideoFiles    `json:"video_files"`
	VideoPictures []VideoPictures `json:"video_pictures"`
}

type PopularVideos struct {
	Page         int32   `json:"page"`
	PerPage      int32   `json:"per_page"`
	TotalResults int32   `json:"total_results"`
	Url          string  `json:"url"`
	Videos       []Video `json:"videos"`
}

type VideoFiles struct {
	Id       int32  `json:"id"`
	Quality  string `json:"quality"`
	FileType string `json:"file_type"`
	Width    int32  `json:"width"`
	Height   int32  `json:"height"`
	Link     string `json:"link"`
}

type VideoPictures struct {
	Id      int32  `json:"id"`
	Picture string `json:"picture"`
	Nr      int32  `json:"nr"`
}

func (c *Client) SearchPhotos(query string, page int32, perPage int32) (*SearchResult, error) {
	url := fmt.Sprintf(PhotoAPI+"/search?query=%s&per_page=%d&page=%d", query, perPage, page)
	fmt.Println(url)
	resp, err := c.RequestDoWithAUth("GET", url)
	if err != nil {
		fmt.Println("Error occurred while authentication")
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result SearchResult
	if err = json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, err
}

func (c *Client) RequestDoWithAUth(method string, url string) (*http.Response, error) {
	// setting the request
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	//adding authorization header
	req.Header.Add("Authorization", c.Token)
	//making the request to the server
	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	times, err := strconv.Atoi(resp.Header.Get("X-RateLimit-Remaining"))
	if err != nil {
		fmt.Println(err)
		return resp, err
	} else {
		c.RemainingTimes = int32(times)
	}
	return resp, nil
}

func (c *Client) CuratedPhotos(perPage, page int) (*CuratedResult, error) {
	url := fmt.Sprintf(PhotoAPI+"/curated?per_page=%d&page=%d", perPage, page)
	resp, err := c.RequestDoWithAUth("GET", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result CuratedResult

	if err = json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, err
}

func (c *Client) GetPhoto(id int32) (*Photo, error) {
	url := fmt.Sprintf(PhotoAPI+"/photos/%d", id)
	resp, err := c.RequestDoWithAUth("GET", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result Photo
	if err = json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) SearchVideo(query string, perPage int, page int) (*VideoSearchResult, error) {
	url := fmt.Sprintf(VideoAPI+"/search?query=%s&per_page=%d&page=%d", query, perPage, page)
	resp, err := c.RequestDoWithAUth("GET", url)
	fmt.Println(url)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result VideoSearchResult
	if err = json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, err
}

func (c *Client) PopularVideo(perPage, page int) (*PopularVideos, error) {
	url := fmt.Sprintf(VideoAPI+"/popular?per_page=%d&page=%d", perPage, page)
	c.RequestDoWithAUth("GET", url)
	resp, err := c.RequestDoWithAUth("GET", url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result PopularVideos
	if err = json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, err
}

func (c *Client) GetRandomVideo() (*Video, error) {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(1001)
	result, err := c.PopularVideo(1, randNum)
	if err == nil && len(result.Videos) == 1 {
		return &result.Videos[0], nil
	}
	return nil, err
}
func (c *Client) GetRemainingRequestInThisMonth() int32 {
	return c.RemainingTimes
}

func (c *Client) GetRandomPhoto() (*Photo, error) {
	rand.Seed(time.Now().Unix())
	randum := rand.Intn(1001)
	result, err := c.CuratedPhotos(1, randum)
	if err == nil && len(result.Photos) == 1 {
		return &result.Photos[0], nil
	}
	return nil, err
}

func main() {
	var Token = os.Getenv("PexelsToken")

	var c = NewClient(Token)
	//result, err := c.SearchVideo("waves", 1, 50)
	//result, err := c.SearchPhotos("waves", 1, 50)
	//result, err := c.GetRandomPhoto()
	result, err := c.GetRandomVideo()
	if err != nil {
		fmt.Errorf("Search error : %v", err)
	}
	//if result.Page == 0 {
	//	fmt.Errorf("Search result is wrong")
	//}
	fmt.Println(result)
}
