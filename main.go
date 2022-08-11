package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

const (
	PhotoAPI = "https://api.pexels.com/v1"
	VideoAPI = "https://api.pexels.com/vi"
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

func (c *Client) SearchPhotos(query string, page int32, perPage int32) (*SearchResult, error) {
	url := fmt.Sprintf(PhotoAPI+"/search?query=%s&per_page=%d&page=%d", query, perPage, page)
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

func main() {
	os.Setenv("PexelsToken", "")
	var Token = os.Getenv("PexelsToken")

	var c = NewClient(Token)
	result, err := c.SearchPhotos("waves", 1, 50)
	if err != nil {
		fmt.Errorf("Search error : %v", err)
	}
	if result.Page == 0 {
		fmt.Errorf("Search result is wrong")
	}
	fmt.Println(result)
}
