// Tools for parsing an org file into a struct
package org

import (
	"errors"
	"fmt"
	"gtdbot/utils"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func GetOrgStatuses() []string {
	return []string{"TODO", "DONE", "CANCELLED", "BLOCKED", "PROGRESS", "WAITING", "TENTATIVE", "DELEGATED"}
}

type OrgDocument struct {
	Filename   string
	Sections   []Section
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
			return section, nil
		}
	}
	return Section{}, errors.New("Section not found")
}

func (o OrgDocument) PrintAll() {
	fmt.Println("Printing all sections from Document: ", o.Filename)
	for _, section := range o.Sections {
		fmt.Println(section.Header())
		for _, item := range section.Items {
			fmt.Println(item.FullLine(section.IndentLevel + 1))
		}
	}
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

func (o OrgDocument) GetFile() *os.File {
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
	utils.InsertLinesInFile(doc.GetFile(), new_lines, section.StartLine+1)
}

type ParseDebugger struct {
	active bool
}

func (pd ParseDebugger) Println(line ...any) {
	if pd.active {
		fmt.Println(line...)
	}
}

func ParseSectionsFromFile(file_name string, serializer BaseOrgSerializer) ([]Section, error) {
	file := GetOrgFile(file_name)
	all_lines, _ := utils.LinesFromReader(file)
	file.Close()

	var sections []Section
	var header string
	start_line := 0
	in_section := false
	print_debugger := ParseDebugger{active: false}

	var items []OrgTODO
	var item_lines []string
	building_item := false

	for i, line := range all_lines {
		print_debugger.Println("line:", line)
		if !strings.HasPrefix(line, "*") {
			if building_item {
				item_lines = append(item_lines, line)
				print_debugger.Println("Building item: ", line)
				continue
			}
			panic("Rogue line: " + line)
		}

		if in_section && strings.HasPrefix(line, "*") && !strings.HasPrefix(line, "**") {
			// Check if we're into the next section at the same indent level as the header
			if building_item {
				building_item = false
				item, serialize_err := serializer.Serialize(item_lines)
				if serialize_err != nil {
					panic("Error serializing item: " + serialize_err.Error())
				}
				items = append(items, item)
				print_debugger.Println("Adding item: ", item.Summary(), item.Details())
			}
			sections = append(sections, Section{
				Description: CleanHeader(header),
				StartLine:   start_line,
				//IndentLevel: strings.Count(header, "*") + 1,
				IndentLevel: 2,
				Items:       items,
			})

			// cleanup, get ready for next section
			items = []OrgTODO{}
			header = ""
			in_section = false
			building_item = false
		}

		if CheckForHeader("TODO", line, "*") || CheckForHeader("DONE", line, "*") {
			in_section = true
			start_line = i
			header = CleanHeader(line)
			print_debugger.Println("Found Section Header: ", header)
			items = []OrgTODO{}
			building_item = false
			item_lines = []string{}
			continue
		}

		if in_section && strings.HasPrefix(line, "**") {
			if building_item {
				building_item = false
				item, serialize_err := serializer.Serialize(item_lines)
				if serialize_err != nil {
					panic("Error serializing item: " + serialize_err.Error())
				}
				items = append(items, item)
				print_debugger.Println("Adding item: ", item.Summary(), item.Details())

				// print_debugger.Println("Starting to build item: ", line)
				// item_lines = []string{line}
				// building_item = true
			}
			print_debugger.Println("Starting to build item: ", line)
			item_lines = []string{line}
			building_item = true
			continue
		}
	}
	// if we're at the end of the file, we need to add the last section
	if in_section {
		sections = append(sections, Section{
			Description: CleanHeader(header),
			StartLine:   start_line,
			IndentLevel: strings.Count(header, "*") + 1,
			Items:       items,
		})
	}

	// if start_line == 0 {
	//	return sections, errors.New("Did not find parsed section.")
	// }
	return sections, nil
}

type OrgItem struct {
	header  string
	details []string
	status  string
	tags    []string
}

func NewOrgItem(header string, details []string, status string, tags []string) OrgItem {
	return OrgItem{
		header,
		details,
		status,
		tags,
	}
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

func (oi OrgItem) Summary() string {
	return oi.header
}

func (oi OrgItem) ID() string {
	return oi.Details()[0]
}

func (oi OrgItem) CheckDone() bool {
	return oi.GetStatus() == "DONE" || oi.GetStatus() == "CANCELLED"
}

func findOrgTags(line string) []string {
	splits := strings.Split(line, ":")
	if len(splits) == 0 {
		return []string{}
	} else {
		return splits[1 : len(splits)-1]
	}

}

func PrintOrgFile(file *os.File) {
	res, _ := ioutil.ReadAll(file)
	fmt.Println(string(res))
}

type OrgSerializer interface {
	Deserialize(item OrgTODO, indent_level int) []string
	Serialize(lines []string) (OrgTODO, error)
}

type BaseOrgSerializer struct{}

// Implement the OrgSerializer interface with our most generic structs / interfaces

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
	tags := findOrgTags(lines[0])
	return OrgItem{header: lines[0], status: status, details: lines[1:], tags: tags}, nil
}

func findOrgStatus(line string) string {
	for _, status := range GetOrgStatuses() {
		if strings.Contains(line, status) {
			return status
		}
	}
	return ""
}
