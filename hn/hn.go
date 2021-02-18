package hn

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// DefaultBaseURL is the default base URL for HN API.
var DefaultBaseURL = "https://hacker-news.firebaseio.com/v0"

// Client is http client for HN.
type Client struct {
	BaseURL string
	MaxItem int64
}

// NewClient creates an instance of Client.
func NewClient() *Client {
	return &Client{
		BaseURL: DefaultBaseURL,
		MaxItem: 0,
	}
}

// GetMaxItem gets MaxItem.
func (c *Client) GetMaxItem() (int64, error) {
	url := fmt.Sprintf("%s/maxitem.json", c.BaseURL)
	b, err := getAndRead(url)
	if err != nil {
		return 0, err
	}

	var mi MaxItem
	err = json.Unmarshal(b, &mi)
	if err != nil {
		return 0, err
	}

	return int64(mi), nil
}

// GetItem gets Item by ID.
func (c *Client) GetItem(id int64) (Item, error) {
	url := fmt.Sprintf("%s/item/%d.json", c.BaseURL, id)
	b, err := getAndRead(url)
	if err != nil {
		return Item{}, err
	}

	var item Item
	err = json.Unmarshal(b, &item)
	if err != nil {
		return Item{}, err
	}

	return item, nil
}

// getAndRead performs http.Get and reads response body.
func getAndRead(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}
