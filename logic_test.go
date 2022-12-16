package main

import "testing"

func TestCardInSection(t *testing.T) {
	card := Card{Id: "123", Title: "Title"}
	items := []LeankitCardOrgLineItem{
		{"Title", "TODO", "abc.com/1"},
		{"Title2", "TODO", "abc.com/2"},
	}
	section := Section{Description: "TODO", StartLine: 5, Items: items}

	if !CheckCardAlreadyInSection(card, section) {
		t.Errorf("Unable to find the card {} in section: ", card.Title, section.Description)
	}

	new_card := Card{Id: "1234", Title: "Title3"}
	if CheckCardAlreadyInSection(new_card, section) {
		t.Errorf("Found card {} in section {} when it shouldn't havethe card", card.Title, section.Description)
	}

}
