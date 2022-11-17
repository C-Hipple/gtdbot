package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const BOARD_ENGINEERING = "30502076986726"
const BOARD_CORE = "30502078266703"
const LANE_NEEDS_REVIEW = "30502078354994"
const LANE_REVIEW_IN_PROGRESS = "30502079469819"
const LANE_REVIEW_NEEDS_CHANGES = "30502078354995"
const LANE_RECENTLY_FINISHED = "30502078276213"
const LANE_CHRIS_DOING_NOW = "30502079234201"

const BASE_HOST = "https://multimediallc.leankit.com/io"

type Card struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type CardResponse struct {
	// Helper class for parsing json response
	Cards []Card `json:"cards"`
}

func getCards(board_id string, lane_id string) {
	client := http.Client{Timeout: 30 * time.Second}
	url := BASE_HOST + "/card?board=" + board_id
	if lane_id != "" {
		url = url + "&lane=" + lane_id
	}
	fmt.Println(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error! Exiting: ", err)
		os.Exit(1)
	}

	req.SetBasicAuth("chrishipple@multimediallc.com", "")

	resp, err2 := client.Do(req)
	if err2 != nil {
		fmt.Println("Error on Do! Exiting: ", err2)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body_bytes, err3 := ioutil.ReadAll(resp.Body)
	if err3 != nil {
		fmt.Println("Error on Reading body: ", err3)
	}

	if false {
		pretty_print(body_bytes)
	}

	var cards CardResponse
	err4 := json.Unmarshal(body_bytes, &cards)
	if err4 != nil {
		fmt.Println("Error on json decode: ", err4)
	} else {
		fmt.Println("We did it!")
	}
	fmt.Println(cards)
}

func pretty_print(body []byte) {
	var pretty_json bytes.Buffer
	err := json.Indent(&pretty_json, body, "", "\t")
	if err != nil {
		fmt.Println("Error prettyifying json", err)
	}
	fmt.Println(string(pretty_json.Bytes()))
}

func main_leankit() {
	getCards(BOARD_CORE, LANE_CHRIS_DOING_NOW)
}
