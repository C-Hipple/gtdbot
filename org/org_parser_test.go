package org

import (
	"reflect"
	"strings"
	"testing"
)

func Test_findOrgTags(t *testing.T) {
	line := "* TODO Example name  :tag1:tag2:"
	tags := findOrgTags(line)
	if !reflect.DeepEqual([]string{"tag1", "tag2"}, tags) {
		t.Fatalf(":tag1:tag2: parsed to %v", tags)
	}

	line2 := "* TODO Example name  :tag1:"
	tags2 := findOrgTags(line2)
	if !reflect.DeepEqual([]string{"tag1"}, tags2) {
		t.Fatalf(":tag1:tag2: parsed to %v", tags2)
	}
}

func Test_findOrgStatus(t *testing.T) {
	line_1 := "* TODO Example name  :tag1:tag2:"
	status_1 := findOrgStatus(line_1)
	if !(status_1 == "TODO") {
		t.Fatalf("Failed.  TODO Should've been found. Found status_1: %v", status_1)
	}

	line_2 := "* DONE Example name  :tag1:tag2:"
	status_2 := findOrgStatus(line_2)
	if !(status_2 == "DONE") {
		t.Fatalf("Failed.  DONE Should've been found. Found status_2: %v", status_2)
	}

}

func Test_Serialize(t *testing.T) {
	raw := `* TODO dev: Serialize PrivateShow Objects into Dataclass		:chaturbate:
15479
https://github.com/multimediallc/chaturbate/pull/15479
Title: dev: Serialize PrivateShow Objects into Dataclass
Author: C-Hipple
`
	parser := BaseOrgSerializer{}
	item, err := parser.Serialize(strings.Split(raw, "\n"))
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	expected := OrgItem{
		header: "* TODO dev: Serialize PrivateShow Objects into Dataclass		:chaturbate:",
		status: "TODO",
		details: []string{
			"15479",
			"https://github.com/multimediallc/chaturbate/pull/15479",
			"Title: dev: Serialize PrivateShow Objects into Dataclass",
			"Author: C-Hipple",
		},
		tags: []string{"chaturbate"},
	}
	if item.GetStatus() != expected.GetStatus() {
		t.Fatalf("Mismatched status %v-%v", item.GetStatus(), expected.GetStatus())
	}

	if item.Summary() != expected.Summary() {
		t.Fatalf("Mismatched summary %v-%v", item.Summary(), expected.Summary())
	}

	if item.ID() != "15479" {
		t.Fatalf("Mismatched ID %v-%v", item.ID(), "15479")
	}

	if item.CheckDone() {
		t.Fatalf("This isn't done!")
	}
}
