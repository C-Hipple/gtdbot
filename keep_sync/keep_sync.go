package main

// package keep_sync

import (
	"encoding/json"
	"net/http"
	"io/ioutil"
)

const BASE_URL = "https://keep.googleapis.com/v1/"
const LIST_URL = BASE_URL + "notes"

type KeepNote struct {
	Name  string `json:"name"`
	Title string `json:"title"`
	Body  Body   `json:"body"`
}

type Body struct {
	Text string     `json:"text"`
	List []ListItem `json":list`
}

type ListItem struct {
	Checked  bool
	Text     string
	Children []ListItem
}

type KeepClient struct {
	client   *http.Client
}

func NewKeepClient() *KeepClient {
	client := http.Client{}

	return &KeepClient{&client}
}

func (kc KeepClient) List() ([]*KeepNote, error) {
	var output []*KeepNote
	res, err := kc.client.Get(kc.base_url)
	if err != nil {
		return output, err
	}

	body_bytes, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(body_bytes, output)
	if err != nil {
		return output, err
	}

	return output, nil
}

func main() {
	client := NewKeepClient()
}
