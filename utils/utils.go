package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	//"io/ioutil"
	//"os"
)

// TODO consider being able to do multiiple operations sorted by at_line_number without opening/closing file
func ReplaceLinesInFile(file *os.File, new_lines []string, at_line_number int) error {
	lines, _ := LinesFromReader(file)
	path, _ := filepath.Abs(file.Name())
	file_content := replaceLines(lines, new_lines, at_line_number)
	return os.WriteFile(path, []byte(file_content), 0644)
}

func replaceLines(existing_lines []string, new_lines []string, at_line_number int) string {
	// Helper so we don't need a file for unit tests
	file_content := ""
	// fmt.Println("Calling replace with the newlines: ", new_lines)
	for i, line := range existing_lines {
		if i == at_line_number {
			file_content += strings.Join(new_lines, "\n")
			continue
		}
		if i >= at_line_number && i < at_line_number + len(new_lines) {
			continue
		}
		file_content += line
		file_content += "\n"
	}
	return file_content
}

func InsertLinesInFile(file *os.File, new_lines []string, at_line_number int) error {
		//new line is at the at_line_number, not after, pushes everything below it.
		lines, _ := LinesFromReader(file)
		file_content := insertLines(lines, new_lines, at_line_number)
		path, _ := filepath.Abs(file.Name())
	return os.WriteFile(path, []byte(file_content), 0644)
}

func insertLines(existing_lines []string, new_lines []string, at_line_number int) string {
	// Helper! for unit tests so we don't need to make a file
	file_content := ""
	for i, line := range existing_lines {
		if i == at_line_number {
			for _, new_line := range new_lines {
				file_content += new_line
				if !strings.HasSuffix(new_line, "\n") {
					file_content += "\n"
				}
			}
		}
		file_content += line
		file_content += "\n"
	}
	return file_content
}

func LinesFromReader(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err := scanner.Err()
	if err != nil {
		return nil, err
	}
	return lines, nil
}

func Contains(check_val string, slice []string) bool {
	for _, val := range slice {
		if check_val == val {
			return true
		}
	}
	return false
}

func PrettyPrint(body []byte) {
	// Pretty print a json bytes array
	var pretty_json bytes.Buffer
	err := json.Indent(&pretty_json, body, "", "\t")
	if err != nil {
		fmt.Println("Error prettyifying json", err)
	}
	fmt.Println(string(pretty_json.Bytes()))
}
