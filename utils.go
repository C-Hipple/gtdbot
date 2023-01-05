package main

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

func ReplaceLineInFile(file *os.File, new_line string, at_line_number int) error {
	lines, _ := LinesFromReader(file)
	file_content := ""
	for i, line := range lines {
		if i == at_line_number {
			file_content += new_line
			continue
		}
		file_content += line
		file_content += "\n"
	}
	path, _ := filepath.Abs(file.Name())
	return os.WriteFile(path, []byte(file_content), 0644)
}

func InsertLinesToFile(file *os.File, new_lines []string, at_line_number int) error {
	//new line is at the at_line_number, not after, pushes everything below it.
	lines, _ := LinesFromReader(file)
	file_content := ""
	for i, line := range lines {
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
	path, _ := filepath.Abs(file.Name())
	return os.WriteFile(path, []byte(file_content), 0644)
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

func pretty_print(body []byte) {
	// Pretty print a json bytes array
	var pretty_json bytes.Buffer
	err := json.Indent(&pretty_json, body, "", "\t")
	if err != nil {
		fmt.Println("Error prettyifying json", err)
	}
	fmt.Println(string(pretty_json.Bytes()))
}
