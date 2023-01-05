package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
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

func (li LeankitCardOrgLineItem) Details() []string {
	return li.Notes
}

func (li LeankitCardOrgLineItem) GetStatus() string {
	return li.Status
}

type OrgTODO interface {
	FullLine(indent_level int) string
	Details() []string
	GetStatus() string
	CheckDone() bool
}

func (li LeankitCardOrgLineItem) CheckDone() bool {
	return li.GetStatus() == "DONE" || li.GetStatus() == "CANCELLED"
}

type Section struct {
	Description string
	StartLine   int
	IndentLevel int // How many ** are we rocking per item (each item, not the header!)
	Items       []OrgTODO
	File        *os.File
}

func (s Section) GetStatus() string {
	if s.CheckAllComplete() {
		return "DONE"
	}
	return "TODO"
}

func (s Section) CalculateDoneRatio() string {
	// This is the [0/0] at the end
	return "[" + strconv.Itoa(len(s.Items)) + "/" + strconv.Itoa(s.DoneCount()) + "]"
}

func (s Section) Header() string {
	header_items := []string{
		strings.Repeat("*", s.IndentLevel-1), // header is 1 less indent than items
		s.GetStatus(),
		s.Description,
		s.CalculateDoneRatio(),
	}

	return strings.Join(header_items, " ")
}

func (s Section) CheckAllComplete() bool {
	return s.DoneCount() == len(s.Items)
}

func (s Section) DoneCount() int {
	var done_count int
	for _, item := range s.Items {
		if item.CheckDone() {
			done_count = done_count + 1
		}
	}
	return done_count
}

func PrintOrgFile(file *os.File) {
	res, _ := ioutil.ReadAll(file)
	fmt.Println(string(res))
}

func ParseOrgFileSection(file *os.File, section_name string, header_indent_level int) (Section, error) {
	org_serializer := OrgSerializer{}
	all_lines, _ := LinesFromReader(file)

	var reviews []LeankitCardOrgLineItem
	start_line := 0
	in_section := false
	building_item := false
	var item_lines []string
	for i, line := range all_lines {
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
			if building_item {
				building_item = false
				item, serialize_err := org_serializer.Serialize(item_lines)
				if serialize_err != nil {
					continue
				}
				reviews = append(reviews, item)
			}
			item_lines = append(item_lines, line)
			building_item = true
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

	OrgFilePath := os.Getenv("OrgFilePath")
	file, err := os.Open(OrgFilePath)
	if err != nil {
		fmt.Println("Error Opening file: ", err)
		os.Exit(1)
	}
	return file
}

func AddTODO(file *os.File, section Section, new_item OrgTODO) {
	serializer := OrgSerializer{}
	// https://siongui.github.io/2017/01/30/go-insert-line-or-string-to-file/#:~:text=If%20you%20want%20to%20insert,the%20end%20of%20the%20string.
	InsertLinesToFile(file, serializer.Deserialize(new_item, section.IndentLevel), section.StartLine+1)
}

func GetOrgSection(section_name string) Section {
	section, err := ParseOrgFileSection(GetOrgFile(), section_name, SectionIndentDepth)
	if err != nil {
		fmt.Println("Error parsing section", section_name, err)
		os.Exit(1)
	}
	return section
}

// func UpdateOrgSectionHeader(section Section) {
// 	var done_count int64
// 		InsertLinesToFile(*section.File, []string{"Header"}, section.)
//       ReplaceLineInFile(*section.File, )

// 	}
// }
