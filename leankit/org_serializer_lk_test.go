package leankit

import (
	"testing"
)

func TestDeserial(t *testing.T) {
	os := OrgLKSerializer{}
	item := Fixture_LeankitCardOrgLineItem(0)
	res := os.Deserialize(item, 2)
	if res[0] != "** TODO "+item.Url+" "+item.Title {
		t.Errorf("Deserialize header is incorrect: %s, %s", res[0], item.FullLine(2))
	}
	details := res[1:]
	target_details := item.Details()
	if len(details) != len(target_details) {
		t.Errorf("Details have different numbers of items: correct: %d, recieved: %d", len(details), len(item.Details()))
	}
	//for i := range len(details) {
	for i := 0; i < len(details); i++ {
		if target_details[i] != details[i] {
			t.Errorf("Element %d of details is wrong.  Should be: %s...Got %s", i, target_details[i], details[i])
		}
	}
}
