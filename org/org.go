package org

import (
	"fmt"
	"os"
	"strings"
)

type OrgTODO interface {
	FullLine(indent_level int) string
	Summary() string
	Details() []string
	GetStatus() string
	CheckDone() bool
	ID() string
	StartLine() int
	LinesCount() int
	Repo() string
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

func CheckTODOInSection(todo OrgTODO, section Section) (int, OrgTODO) {
	// returns the line number if it's found, otherwise returns -1
	serializer := BaseOrgSerializer{}
	at_line := section.StartLine + 1 // account for the section title
	for _, line_item := range section.Items {
		if strings.Contains(line_item.Summary(), todo.Summary()) {
			return at_line, line_item
		}
		if line_item.Summary() == todo.Summary() {
			return at_line, line_item
		}
		at_line += len(serializer.Deserialize(line_item, section.IndentLevel))
	}
	return -1, nil
}
