package main

import (
	"bufio"
	"io"
	//"io/ioutil"
	//"os"
)

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
