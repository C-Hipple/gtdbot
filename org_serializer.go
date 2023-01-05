package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/go-github/v48/github"
)

// Service struct for
type OrgSerializer struct{}

func (os OrgSerializer) Deserialize(item OrgTODO, indent_level int) []string {
	var result []string
	result = append(result, item.FullLine(indent_level))
	result = append(result, item.Details()...)
	fmt.Println(result)
	return result
}

func (os OrgSerializer) Serialize(lines []string) (LeankitCardOrgLineItem, error) {
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

func DeserializePRBody(body *string) []string {
	result := []string{}
	return result
}

func (os OrgSerializer) DeserializePR(pr *github.PullRequest) []string {
	result := []string{}
	result = append(result, *pr.Body)
	return result

}

func (os OrgSerializer) SerializePRToID(pr *github.PullRequest) int {
	// Only does to the ID, don't want to make un-necessary API calls
	return 0
}
