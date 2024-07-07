package utils

import (
	"strings"
	"testing"
)

func Test_ReplaceLines(t *testing.T) {
	existing_lines := strings.Split(`* TODO Code Review
** TODO PR #1 :draft:
the_url
author`, "\n")
	new_lines := []string{"** TODO PR #1", "updated_url"}
    updated := replaceLines(existing_lines, new_lines, 1)
	target :=`* TODO Code Review
** TODO PR #1
updated_url
author
`  // newline at the end
	if updated != target {
		t.Fatalf("Updated lines do not match.  Target: %v \nActual %v", target, updated)
	}
}


func Test_InsertLines(t *testing.T) {
	existing_lines := strings.Split(`* TODO Code Review
** TODO PR #1 :draft:
the_url
author`, "\n")
	new_lines := []string{"sub-header", "add. sub-header"}
    updated := insertLines(existing_lines, new_lines, 2)
	target :=`* TODO Code Review
** TODO PR #1 :draft:
sub-header
add. sub-header
the_url
author
`  // newline at the end
	if updated != target {
		t.Fatalf("Updated lines do not match.  Target: %v \nActual %v", target, updated)
	}
}
