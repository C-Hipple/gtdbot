package utils

import (
	"fmt"
	"strings"
	"testing"
)

func Test_ReplaceLines(t *testing.T) {
	existing_lines := []string{"* TODO Code Review",
		"** TODO PR #1",
		"the_url",
		"author",
		"misc_info",
		"\n"}
	new_lines := []string{"** TODO PR #1", "updated_url"}
	updated := replaceLines(existing_lines, new_lines, 1)
	target := []string{"* TODO Code Review",
		"** TODO PR #1",
		"updated_url",
		"author",
		"misc_info",
		"\n"}
	for _, ele := range Zip(updated, target) {
		fmt.Println(ele.First, ele.Second)
		if ele.First != ele.Second {
			fmt.Printf("len target: %d, len actual: %d", len(target), len(updated))
			t.Fatalf("Updated lines do not match.  Target: \n%v Actual \n%v", target, updated)
		}

	}
}

func Test_ReplaceLinesInMiddle(t *testing.T) {
	existing_lines := []string{"* TODO Code Review",
		"** TODO PR #1",
		"the_url",
		"author",
		"misc_info",
		"\n"}
	new_lines := []string{"updated_url", "updated_author"}
	updated := replaceLines(existing_lines, new_lines, 2)

	target := []string{"* TODO Code Review",
		"** TODO PR #1",
		"updated_url",
		"updated_author",
		"misc_info",
		"\n"}
	for _, ele := range Zip(updated, target) {
		fmt.Println(ele.First, ele.Second)
		if ele.First != ele.Second {
			fmt.Printf("len target: %d, len actual: %d", len(target), len(updated))
			t.Fatalf("Updated lines do not match.  Target: \n%v Actual \n%v", target, updated)
		}

	}
}

func Test_InsertLines(t *testing.T) {
	existing_lines := strings.Split(`* TODO Code Review
** TODO PR #1 :draft:
the_url
author`, "\n")
	new_lines := []string{"sub-header", "add. sub-header"}
	updated := insertLines(existing_lines, new_lines, 2)
	target := `* TODO Code Review
** TODO PR #1 :draft:
sub-header
add. sub-header
the_url
author
` // newline at the end
	if updated != target {
		t.Fatalf("Updated lines do not match.  Target: %v \nActual %v", target, updated)
	}
}

func Test_RemoveLines(t *testing.T) {
	existing_lines := []string{
		"* TODO PR #1 <title> :draft:",
		"sub-header",
		"add. sub-header",
		"url",
		"the_author",
		"*** Body",
		"",
		"In this PR We do things",
	}

	removed := removeLines(existing_lines, 0, 1)
	if removed[0] != "sub-header" {
		t.Fatal("Failed to remove the first line")
	}
	if len(removed) != len(existing_lines)-1 {

		t.Fatal("Removed too many lines")

	}

	removed = removeLines(existing_lines, 0, 0)
	if removed[0] != existing_lines[0] {
		t.Fatal("Remove 0 lines failed")
	}

	if len(removed) != len(existing_lines) {
		t.Fatal("Removed lines when shouldn't ahve")
	}
}
