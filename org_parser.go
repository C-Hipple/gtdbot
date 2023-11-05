// Tools for parsing an org file into a struct
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"io/ioutil"
	"errors"
)


func GetOrgStatuses () []string {
	return []string{"TODO", "DONE", "CANCELLED", "BLOCKED", "PROGRESS", "WAITING", "TENTATIVE", "DELEGATED"}
}


type OrgDocument struct {
	Filename string
	Sections []Section
	Serializer OrgSerializer
}

func GetOrgDocument(file_name string) OrgDocument {
	serializer := BaseOrgSerializer{}
	sections, err := ParseSectionsFromFile(file_name, serializer)
	if err != nil {
		fmt.Println("Error parsing sections from file: ", err)
		os.Exit(1)
	}
	doc := OrgDocument{Filename: file_name, Sections: sections, Serializer: serializer}
	return doc
}

func (o OrgDocument) Refresh() {
	serializer := BaseOrgSerializer{}
	sections, err := ParseSectionsFromFile(o.Filename, serializer)
	if err != nil {
		fmt.Println("Error parsing sections from file: ", err)
		os.Exit(1)
	}
	o.Sections = sections
}

func (o OrgDocument) GetSection(section_name string) (Section, error) {
	for _, section := range o.Sections {
		if section.Description == section_name {
			fmt.Println("Found section: ", section)
			return section, nil
		}
	}
	return Section{}, errors.New("Section not found")
}

type Section struct {
	Description string
	StartLine   int
	IndentLevel int // How many ** are we rocking per item (each item, not the header!)
	Items       []OrgTODO
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

func (o OrgDocument) GetFile () *os.File {
	return GetOrgFile(o.Filename)
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

func AddTODO(doc OrgDocument, section Section, new_item OrgTODO) {
	section.Items = append(section.Items, new_item)
	new_lines := doc.Serializer.Deserialize(new_item, section.IndentLevel)
	InsertLinesToFile(doc.GetFile(), new_lines, section.StartLine+1)
}


func ParseSectionsFromFile(file_name string, serializer BaseOrgSerializer) ([]Section, error) {
	file := GetOrgFile(file_name)
	all_lines, _ := LinesFromReader(file)
	file.Close()

	var sections []Section
	var header string
	start_line := 0
	in_section := false

	var items []OrgTODO
	var item_lines []string
	building_item := false

	for i, line := range all_lines {
		if !strings.HasPrefix(line, "*") {
			if building_item {
				item_lines = append(item_lines, line)
			}
		}

		if in_section && strings.HasPrefix(line, "*") && !strings.HasPrefix(line, "**") {
			// Check if we're into the next section at the same indent level as the header
			sections = append(sections, Section {
				Description: CleanHeader(header),
					StartLine: start_line,
					IndentLevel: strings.Count(header, "*") + 1,
					Items: items,
				})

			// cleanup, get ready for next section
			items = []OrgTODO{}
			header = ""
			in_section = false
		}

		if CheckForHeader("TODO", line, "*") {
			in_section = true
			start_line = i
			header = line
			continue
		}

		if in_section && strings.HasPrefix(line, "**") {
			if building_item {
				building_item = false
				item, serialize_err := serializer.Serialize(item_lines)
				if serialize_err != nil {
					continue
				}
				items = append(items, item)
			}
			item_lines = append(item_lines, line)
			building_item = true
		}
	}
	// if start_line == 0 {
	//	return sections, errors.New("Did not find parsed section.")
	// }
	return sections, nil
}

type OrgItem struct {
	header string
	details  []string
	status   string
}

// Implement the OrgTODO Interface for OrgItem
func (oi OrgItem) FullLine(indent_level int) string {
	return strings.Repeat("*", indent_level) + oi.header
}

func (oi OrgItem) Details() []string {
	return oi.details
}

func (oi OrgItem) GetStatus() string {
	return oi.status
}

func (oi OrgItem) CheckDone() bool {
	return oi.GetStatus() == "DONE" || oi.GetStatus() == "CANCELLED"
}



func PrintOrgFile(file *os.File) {
	res, _ := ioutil.ReadAll(file)
	fmt.Println(string(res))
}

type OrgSerializer interface {
	Deserialize(item OrgTODO, indent_level int) []string
	Serialize(lines []string) (OrgTODO, error)
}

type BaseOrgSerializer struct {}

func (bos BaseOrgSerializer) String() string {
	return "BaseOrgSerializer"
}

func (bos BaseOrgSerializer) Deserialize(item OrgTODO, indent_level int) []string {
	var result []string
	result = append(result, item.FullLine(indent_level))
	result = append(result, item.Details()...)
	return result
}

func (bos BaseOrgSerializer) Serialize(lines []string) (OrgTODO, error) {
	// each one has the format ** TODO URL Title.  Check stars to allow for auxillary text between items
	if len(lines) == 0 {
		return OrgItem{}, errors.New("No Lines passed for serialization")
	}
	status := findOrgStatus(lines[0])
	return OrgItem{header: lines[0], status: status, details: lines[1:]}, nil
}


func findOrgStatus(line string) string {
	for _, status := range GetOrgStatuses() {
		if strings.Contains(line, "TODO") {
			return status
		}
	}
	return ""

}

func ParseOrgFileSection(file *os.File, section_name string, header_indent_level int) (Section, error) {
	org_serializer := OrgLKSerializer{}
	all_lines, _ := LinesFromReader(file)

	var found_items []LeankitCardOrgLineItem
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
				found_items = append(found_items, item)
			}
			item_lines = append(item_lines, line)
			building_item = true
		}
	}
	if start_line == 0 {
		return Section{}, errors.New("Did not find parsed section.")

	}
	sec := Section{Description: section_name, IndentLevel: 3, StartLine: start_line}
	// would've thought I could do Items: found_items ^^ but it's a typing issue :(
	for _, review := range found_items {
		sec.Items = append(sec.Items, review)
	}
	return sec, nil
}
