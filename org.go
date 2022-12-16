package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const OrgFilePath = "/Users/chrishipple/gtd/gtd.org"
const SectionIndentDepth = 2

type LeankitCardOrgLineItem struct {
	Title  string
	Status string
	Url    string
	// Notes []string will be the other things below it..
}

func (cr LeankitCardOrgLineItem) FullLine(level int) string {
	return strings.Repeat("*", level) + " " + cr.Status + " " + cr.Url + " " + cr.Title
}

type OrgTODO interface {
	FullLine(int) string
}

type Section struct {
	Description string
	StartLine   int
	IndentLevel int // How many ** are we rocking per item (each item, not the header!)
	Items       []OrgTODO
	File        *os.File
}

func PrintOrgFile(file *os.File) {
	res, _ := ioutil.ReadAll(file)
	fmt.Println(string(res))
}

func ParseOrgFileSection(file *os.File, section_name string, header_indent_level int) (Section, error) {
	split, _ := LinesFromReader(file)
	in_section := false
	var reviews []LeankitCardOrgLineItem
	start_line := 0

	for i, line := range split {
		// fmt.Println(line)
		if !strings.HasPrefix(line, "*") {
			continue // this is helper text or some other nonsense
		}
		stars := strings.Repeat("*", header_indent_level)
		if in_section && strings.HasPrefix(line, stars) && !strings.HasPrefix(line, stars+"*") {
			// Check if we're into the next section at the same indent level as the header
			in_section = false
			break
		}
		if CheckForHeader(section_name, line, stars) {
			in_section = true
			start_line = i
			continue
		}
		if in_section && strings.HasPrefix(line, stars+"*") {
			// each one has the format ** TODO URL Title.  Check stars to allow for auxillary text between items
			split_line_item := strings.Split(line, " ")
			if len(split_line_item) < 4 {
				continue // This is not from a leankit card from this bot, can be ignored.
			}
			reviews = append(reviews,
				LeankitCardOrgLineItem{split_line_item[3], split_line_item[1], split_line_item[2]})
		}
	}
	if start_line == 0 {
		return Section{}, errors.New("Did not find parsed section.")

	}
	sec := Section{Description: section_name, IndentLevel: 3, StartLine: start_line}
	// would've thought I could do Items: reviews ^^ but it's a typing issue :(
	for _, review := range reviews {
		sec.Items = append(sec.Items, review)
	}
	return sec, nil
}

func CheckForHeader(section_name string, line string, stars string) bool {
	prefix := strings.HasPrefix(line, stars+" TODO ") || strings.HasPrefix(line, stars+" DONE ")
	return prefix && strings.Contains(line, section_name)
}

func GetOrgFile() *os.File {
	file, err := os.Open(OrgFilePath)
	if err != nil {
		fmt.Println("Error Opening file: ", err)
		os.Exit(1)
	}
	return file
}

func AddTODO(file *os.File, section Section, new_item OrgTODO) {
	// https://siongui.github.io/2017/01/30/go-insert-line-or-string-to-file/#:~:text=If%20you%20want%20to%20insert,the%20end%20of%20the%20string.
	InsertLineToFile(file, new_item.FullLine(section.IndentLevel), section.StartLine+1)
}

func GetOrgSection(section_name string) Section {
	section, err := ParseOrgFileSection(GetOrgFile(), section_name, SectionIndentDepth)
	if err != nil {
		fmt.Println("Error parsing section", section_name, err)
		os.Exit(1)
	}
	return section
}

func InsertLineToFile(file *os.File, newline string, at_line_number int) error {
	//new line is at the at_line_number, not after, pushes everything below it.
	lines, _ := LinesFromReader(file)
	file_content := ""
	for i, line := range lines {
		if i == at_line_number {
			file_content += newline
			if !strings.HasSuffix(newline, "\n") {
				file_content += "\n"
			}
		}
		file_content += line
		file_content += "\n"
	}
	path, _ := filepath.Abs(file.Name())
	return os.WriteFile(path, []byte(file_content), 0644)
}
