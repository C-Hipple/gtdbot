package leankit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	dy
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const BOARD_ENGINEERING = "30502076986726"
const BOARD_CORE = "30502078266703"
const LANE_NEEDS_REVIEW = "30502078354994"
const LANE_REVIEW_IN_PROGRESS = "30502079469819"
const LANE_REVIEW_NEEDS_CHANGES = "30502078354995"
const LANE_RECENTLY_FINISHED = "30502078276213"
const LANE_CHRIS_DOING_NOW = "30502079234201"

// Engineering Board Leankit Lanes
const LANE_ENG_CORE_REVIEW = "30502076989618"
const LANE_ENG_CORE_NEEDS_REVIEW = "30502076991756"

const MY_USER_ID = "30502079267931"

var CORE_ACTIVE_LANES = [...]string{LANE_CHRIS_DOING_NOW, LANE_NEEDS_REVIEW, LANE_REVIEW_IN_PROGRESS, LANE_REVIEW_NEEDS_CHANGES}

// go doesn't do constant slices, makes sense
//const CORE_REVIEW_LANES = []string{"30502076991756","30502078267297","30502078272818","30502077307016","30502078267298",}

const BASE_HOST = "https://owner.leankit.com"
const API_BASE = BASE_HOST + "/io"

type Card struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Links []Link `json:"externalLinks"`
	Users []User `json:"assignedUsers"`
	Body  string `json:"description"`
	Board Board  `json:"board"`
}

type User struct {
	Id       string `json:"id"`
	FullName string `json:"fullName"`
}

type Link struct {
	Label string `json:"label"`
	Url   string `json:"url"`
}

type Board struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

func (c Card) FullLine(indent_level int) string {
	return strings.Repeat("*", indent_level) + " TODO " + c.URL() + " " + c.Title
}

func (c Card) GetCleanBody() string {
	// Removes <p> and </p> html tags that LK api puts on description field.
	return c.Body[3 : len(c.Body)-4]
}

func (c Card) Details() []string {
	var details []string
	if len(c.Users) > 0 {
		details = append(details, "Assigned User: "+c.Users[0].FullName)
	}
	if len(c.Links) > 0 {
		details = append(details, "PR: "+c.Links[0].Url)
	}
	if len(c.Body) > 0 {
		details = append(details, "Card Description: \n")
		details = append(details, c.GetCleanBody())
	}
	return details
}

func (c Card) URL() string {
	return BASE_HOST + "/card/" + c.Id
}

func (c Card) UserID() string {
	if len(c.Users) == 0 {
		return ""
	}
	return c.Users[0].Id
}

func (c Card) GetCardChildren() []Card {
	url := API_BASE + "/card/" + c.Id + "/connection/children"
	fmt.Println(url)
	children, err := MakeCardsAPICall(url, []Filter{})
	if err != nil {
		// better handling?  children is already empty card slice
		return []Card{}
	}
	return children
}

func (c Card) GetStatus() string {
	return "TODO"
}

func (c Card) CheckDone() bool {
	return false
}

func (c Card) IsEngineeringCard() bool {
	return c.Board.Id == BOARD_ENGINEERING
}

func (l Link) PRNumber() int {
	splits := strings.Split(l.Url, "/")
	number, err := strconv.Atoi(splits[len(splits)-1])
	if err != nil {
		fmt.Println("Error parsing PR Number from URL:", l.Url)
		return 0
	}
	return number
}

type CardResponse struct {
	// Helper struct for deserializing json response
	Cards []Card `json:"cards"`
}

func SerializeResponseToCards(resp *http.Response) []Card {
	body_bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error on Reading body: ", err)
	}

	if false {
		pretty_print(body_bytes)
	}

	var cards CardResponse
	err = json.Unmarshal(body_bytes, &cards)
	if err != nil {
		fmt.Println("Error on json decode: ", err)
		return []Card{}
	}

	return cards.Cards
}

type Filter func([]Card) []Card

func getCards(board_id string, lane_ids []string, filters []Filter) []Card {
	url := API_BASE + "/card?board=" + board_id

	if len(lane_ids) > 0 {
		if len(lane_ids) > 1 {
			url = url + "&lanes="
		} else {
			url = url + "&lane="
		}
		url = url + strings.Join(lane_ids, ",")
	}
	//fmt.Println(url)
	cards, err := MakeCardsAPICall(url, filters)
	if err != nil {
		fmt.Println("Error!")
	}
	return cards
}

func MakeCardsAPICall(url string, filters []Filter) ([]Card, error) {
	client := http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error! Exiting: ", err)
		os.Exit(1)
	}
	username := os.Getenv("GTDBOT_LK_API_USERNAME")
	password := os.Getenv("GTDBOT_LK_API_PASSWORD")
	if !validate_auth(username, password) {
		fmt.Println("Username and Password are not set!")
		os.Exit(1)
	}
	req.SetBasicAuth(username, password)

	resp, err2 := client.Do(req)
	if err2 != nil || resp.StatusCode > 200 {
		txt := fmt.Sprintf("Status Code: %d Means %s", resp.StatusCode, http.StatusText(resp.StatusCode))
		fmt.Println(txt)
		return []Card{}, errors.New("Did not get a 200 status code from API: " + txt)
	}
	defer resp.Body.Close()

	cards := SerializeResponseToCards(resp)
	for _, filter_func := range filters {
		cards = filter_func(cards)
	}
	return cards, nil
}

// todo use partial function
func MyUserFilter(cards []Card) []Card {
	CardsOut := []Card{}
	for _, card := range cards {
		if card.UserID() == MY_USER_ID {
			CardsOut = append(CardsOut, card)
		}
	}

	return CardsOut
}

func NotMeFilter(cards []Card) []Card {
	CardsOut := []Card{}
	for _, card := range cards {
		fmt.Println(card.Title, card.UserID(), MY_USER_ID)
		if card.UserID() != MY_USER_ID {
			CardsOut = append(CardsOut, card)
		}
	}
	return CardsOut
}

func lk_main() {
	cards := getCards(BOARD_CORE, []string{LANE_RECENTLY_FINISHED}, []Filter{MyUserFilter})
	for _, card := range cards {
		fmt.Println(card)
		if len(card.Links) > 0 {
			fmt.Println(card.Links[0].PRNumber())
		}
	}

}

func CardFromOrgLineItem(line_item LeankitCardOrgLineItem) Card {
	fmt.Println("Card From Org Line Item Not Implemented! ")
	os.Exit(1)
	return Card{}
}

func validate_auth(username string, password string) bool {
	return !(username == "" && password == "")
}
