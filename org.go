package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

const SectionIndentDepth = 2

type LeankitCardOrgLineItem struct {
	Title  string
	Status string
	Url    string
	Notes  []string // will be the other things below it..
}

func (li LeankitCardOrgLineItem) ID() string {
	id_regex, _ := regexp.Compile("[0-9]+")
	return id_regex.FindString(li.Url)
}

func (li LeankitCardOrgLineItem) FullLine(indent_level int) string {
	return strings.Repeat("*", indent_level) + " " + li.Status + " " + li.Url + " " + li.Title
}

func (li LeankitCardOrgLineItem) Summary() string {
	return li.Title
}

func (li LeankitCardOrgLineItem) Details() []string {
	return li.Notes
}

func (li LeankitCardOrgLineItem) GetStatus() string {
	return li.Status
}

type OrgTODO interface {
	FullLine(indent_level int) string
	Summary() string
	Details() []string
	GetStatus() string
	CheckDone() bool
	ID() string
}

func InterfaceCheck(a OrgTODO) bool {
	return true
}

func useInterface(a OrgItem) bool {
	return InterfaceCheck(a)
}

func (li LeankitCardOrgLineItem) CheckDone() bool {
	return li.GetStatus() == "DONE" || li.GetStatus() == "CANCELLED"
}

func CleanHeader(line string) string {
	line = strings.ReplaceAll(line, "*", "")
	line = strings.ReplaceAll(line, "TODO", "")
	line = strings.ReplaceAll(line, "DONE", "")
	line = strings.TrimSpace(line)
	if strings.Contains(line, "[") {
		line = strings.Split(line, "[")[0]
	}
	return line
}

func CheckForHeader(section_name string, line string, stars string) bool {
	prefix := strings.HasPrefix(line, stars+" TODO ") || strings.HasPrefix(line, stars+" DONE ")
	return prefix && strings.Contains(line, section_name)
}

func GetOrgFile(filename string) *os.File {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory: ", err)
		os.Exit(1)
	}
	org_file_path := home + "/gtd/" + filename

	file, err := os.Open(org_file_path)
	if err != nil {
		fmt.Println("Error Opening Org filefile: ", file, err)
		os.Exit(1)
	}
	return file
}
