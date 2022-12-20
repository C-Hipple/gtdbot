package main

import (
	"errors"
	"fmt"
	"strings"
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

func (os OrgSerializer) Serialize(line string) (LeankitCardOrgLineItem, error) {
	// each one has the format ** TODO URL Title.  Check stars to allow for auxillary text between items
	split_line_item := strings.Split(line, " ")
	if len(split_line_item) < 4 {
		// This is not from a leankit card from this bot, can be ignored.
		return LeankitCardOrgLineItem{}, errors.New("Cannot Serialize line: " + line)
	}
	return LeankitCardOrgLineItem{split_line_item[3], split_line_item[1], split_line_item[2]}, nil
}
