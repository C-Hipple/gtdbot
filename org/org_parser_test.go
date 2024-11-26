package org

import (
	"fmt"
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
	raw := `** TODO dev: Serialize PrivateShow Objects into Dataclass		:chaturbate:
15479
https://github.com/multimediallc/chaturbate/pull/15479
Title: dev: Serialize PrivateShow Objects into Dataclass
Author: C-Hipple
*** BODY
abc

def

ge
`

	parser := BaseOrgSerializer{}
	item, err := parser.Serialize(strings.Split(raw, "\n"), 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	expected := OrgItem{
		header: "** TODO dev: Serialize PrivateShow Objects into Dataclass		:chaturbate:",
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
	if len(item.Details()) != 11 {
		t.Fatalf("Wrong length of details! %v", len(item.Details()))
	}
}

func Test_LineNumbers(t *testing.T) {
	raw := `** TODO dev: Serialize PrivateShow Objects into Dataclass		:chaturbate:
15479
https://github.com/multimediallc/chaturbate/pull/15479
Title: dev: Serialize PrivateShow Objects into Dataclass
Author: C-Hipple
*** BODY
abc

def

ge
`
	// doc :=
	parser := BaseOrgSerializer{}
	item, err := parser.Serialize(strings.Split(raw, "\n"), 0)
	if err != nil {
		t.Fatalf("Error on Serialized: %s", err)
	}
	fmt.Println("Details: ")
	fmt.Println(item.Details())
	fmt.Println("Done ")
}

func Test_ParseSections(t *testing.T) {
	raw := `* TODO Team Reviews [0/2]
** TODO dev: PR 1 :chaturbate:
15479
https://github.com/multimediallc/chaturbate/pull/15479
Title: dev: PR-1
Author: C-Hipple
*** BODY
abc

def

ge
** TODO feature: PR 2 :chaturbate:
15480
https://github.com/multimediallc/chaturbate/pull/15480
Title: feature: PR-2
Author: C-Hipple
*** BODY
abc2

def2

ge2
open line end

* TODO My Review Requests [0/2]
** TODO dev: PR 3 :chaturbate:
15479
https://github.com/multimediallc/chaturbate/pull/15479
Title: dev: PR-3
Author: C-Hipple
*** BODY
abc

def
** TODO feature: PR 4 :chaturbate:
15480
https://github.com/multimediallc/chaturbate/pull/15480
Title: feature: PR-4
Author: C-Hipple
*** BODY
short body
`
	// fmt.Println(raw)
	raw_lines := strings.Split(raw, "\n")
	sections, err := ParseSectionsFromLines(raw_lines, BaseOrgSerializer{})
	if err != nil {
		t.Fatalf("Error parsing sections %v", err)
	}

	if len(sections) != 2 {
		t.Fatalf("Wrong length of sections ! %v", len(sections))
	}

	section_team_review := sections[0]
	section_my_review := sections[1]

	if len(section_my_review.Items) != 2 {
		t.Fatalf("Wrong length of my review items! %v", len(section_my_review.Items))
	}

	if len(section_team_review.Items) != 2 {
		t.Fatalf("Wrong length of team review items! %v", len(section_team_review.Items))
	}

	if section_team_review.StartLine != 0 {
		t.Fatalf("Wrong start line of first section %v", section_team_review.StartLine)
	}

	if section_my_review.StartLine != 25 {
		t.Fatalf("Wrong start line of second section %v", section_my_review.StartLine)
	}

	// Test item start lines
	if section_team_review.Items[0].StartLine() != 1 {
		t.Fatalf("Wrong start line of first item in first section %v", section_team_review.Items[0].StartLine() )
	}

	if section_team_review.Items[1].StartLine() != 12 {
		t.Fatalf("Wrong start line of second item in first section %v", section_team_review.Items[1].StartLine() )
	}

	if section_my_review.Items[0].StartLine() != 26 {
		t.Fatalf("Wrong start line of first item in second section %v", section_my_review.Items[0].StartLine() )
	}

	if section_my_review.Items[1].StartLine() != 35 {
		t.Fatalf("Wrong start line of second item in second section %v", section_my_review.Items[1].StartLine() )
	}
}
