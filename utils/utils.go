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
	updated_lines := strings.Join(FixNewLineEndings(replaceLines(lines, new_lines, at_line_number)), "")
	return os.WriteFile(path, []byte(updated_lines), 0644)
}

func FixNewLineEndings(lines []string) []string {
	// I'm sloppy about what does or does not have a \n at the end
	// and thsi function paves over that
	out_lines := []string{}
	for _, line := range lines {
		if strings.Contains(line, "\n") {
			out_lines = append(out_lines, line)
		} else {
			out_lines = append(out_lines, line+"\n")
		}
	}
	return out_lines
}

func replaceLines(existing_lines []string, new_lines []string, at_line_number int) []string {
	var out_lines []string
	j := 0
	for i, line := range existing_lines {
		if i < at_line_number {
			out_lines = append(out_lines, line)
		}
		if i >= at_line_number && i < at_line_number+len(new_lines) {
			out_lines = append(out_lines, new_lines[j])
			j += 1
		}
		if i >= at_line_number+len(new_lines) {
			out_lines = append(out_lines, line)
		}
	}
	return out_lines
}

func InsertLinesInFile(file *os.File, new_lines []string, at_line_number int) (int, error) {
	// If at_line_number is -1, do it at the end
	// new line is at the at_line_number, not after, pushes everything below it.
	// Returns at_line_number back to user so we know if we added it at the end.
	lines, _ := LinesFromReader(file)
	if at_line_number == -1 {
		at_line_number = len(lines)
	}
	file_content := insertLines(lines, new_lines, at_line_number)
	path, _ := filepath.Abs(file.Name())
	err := os.WriteFile(path, []byte(file_content), 0644)
	if err != nil {
		return 0, err
	}
	return at_line_number, nil
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

func RemoveLinesInFile(file *os.File, at_line_number int, remove_count int) error {
	lines, _ := LinesFromReader(file)
	path, _ := filepath.Abs(file.Name())
	updated_lines := strings.Join(FixNewLineEndings(removeLines(lines, at_line_number, remove_count)), "")
	return os.WriteFile(path, []byte(updated_lines), 0644)
}

func removeLines(existing_lines []string, at_line_number int, remove_count int) []string {
	output_lines := []string{}
	for i, line := range existing_lines {
		if i < at_line_number {
			output_lines = append(output_lines, line)
		} else if i >= at_line_number && i < at_line_number+remove_count {
			continue
		} else {
			output_lines = append(output_lines, line)
		}
	}
	return output_lines
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

type Pair[T, U any] struct {
	First  T
	Second U
}

func Zip[T, U any](ts []T, us []U) []Pair[T, U] {
	if len(ts) != len(us) {
		// TODO: consider handling if different lengths
		panic("slices have different length")
	}
	pairs := make([]Pair[T, U], len(ts))
	for i := 0; i < len(ts); i++ {
		pairs[i] = Pair[T, U]{ts[i], us[i]}
	}
	return pairs
}

func Map[T, U any](ts []T, fn func(T) U) []U {
	result := make([]U, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
}
