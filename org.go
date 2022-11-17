package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const orgFilePath = "/users/chrishipple/gtd/gtd.org" //TODO: expand relative path

type CodeReview struct {
	Title  string
	Status string
	Url    string
}

func Hello() string {
	return orgFilePath
}

func PrintOrgFile(file *os.File) {
	res, _ := ioutil.ReadAll(file)
	fmt.Println(string(res))
}

func ParseCodeReview(file *os.File) []CodeReview {
	res, _ := ioutil.ReadAll(file)
	split := strings.Split(string(res), "\n")
	inCRSection := false
	var reviews []CodeReview

	for _, str := range split {
		if !strings.HasPrefix(str, "*") {
			continue // this is helper text or some other nonsense
		}
		if inCRSection && strings.HasPrefix(str, "**") && !strings.HasPrefix(str, "***") {
			inCRSection = false
			break
		}
		if strings.HasPrefix(str, "** TODO ") && strings.Contains(str, "Code Review") {
			inCRSection = true
			continue
		}
		if inCRSection {
			// each one has the format ** TODO URL Title
			split_line_item := strings.Split(str, " ")
			reviews = append(reviews,
				CodeReview{split_line_item[3], split_line_item[1], split_line_item[2]})
		}
	}
	return reviews
}

func GetOrgFile() *os.File {
	file, err := os.Open(orgFilePath)
	if err != nil {
		fmt.Println("Error Opening file: ", err)
		os.Exit(1)
	}
	return file
}

// https://siongui.github.io/2017/01/30/go-insert-line-or-string-to-file/#:~:text=If%20you%20want%20to%20insert,the%20end%20of%20the%20string.


