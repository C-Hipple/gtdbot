package main

import "testing"

func TestCardInSection(t *testing.T) {
	section := Fixture_Section()
	card := Fixture_Card("1")

	if !CheckCardAlreadyInSection(card, section) {
		t.Errorf("Unable to find the card in section: %s - %s", card.Title, section.Description)
	}

	new_card := Card{Id: "1234", Title: "Title3"}
	if CheckCardAlreadyInSection(new_card, section) {
		t.Errorf("Found card %s in section %s when it shouldn't havethe card", card.Title, section.Description)
	}
}
