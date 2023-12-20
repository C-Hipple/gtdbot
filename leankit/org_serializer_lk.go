package leankit

import (
	"errors"
	"fmt"
	"strings"
)

// Service struct for
type OrgLKSerializer struct{}

func (os OrgLKSerializer) Deserialize(item OrgTODO, indent_level int) []string {
	var result []string
	result = append(result, item.FullLine(indent_level))
	result = append(result, item.Details()...)
	fmt.Println(result)
	return result
}

func (os OrgLKSerializer) Serialize(lines []string) (LeankitCardOrgLineItem, error) {
	// each one has the format ** TODO URL Title.  Check stars to allow for auxillary text between items
	if len(lines) == 0 {
		return LeankitCardOrgLineItem{}, errors.New("No Lines passed for serialization")
	}
	header_line := lines[0]
	notes := lines[1:]

	split_line_item := strings.Split(header_line, " ")
	if len(split_line_item) < 4 {
		// This is not from a leankit card from this bot, can be ignored.
		return LeankitCardOrgLineItem{}, errors.New("Cannot Serialize header_line: " + header_line)
	}
	title := split_line_item[3]
	status := split_line_item[1]
	url := split_line_item[2]

	return LeankitCardOrgLineItem{title, status, url, notes}, nil
}
