package input

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func CsvStrings(path string) (data []string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return data, err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ",")
		data = append(data, line...)
	}
	return data, nil
}

func CsvInts(path string) (data []int, err error) {
	file, err := os.Open(path)
	if err != nil {
		return data, err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ",")
		for _, word := range line {
			num, err := strconv.Atoi(word)
			if err == nil {
				data = append(data, num)
			}
		}
	}
	return data, nil
}
