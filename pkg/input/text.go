package input

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"strings"
)

func splitOnDelimiter(delim string) func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	searchBytes := []byte(delim)
	searchLen := len(searchBytes)
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		dataLen := len(data)

		// Return nothing if at end of file and no data passed
		if atEOF && dataLen == 0 {
			return 0, nil, nil
		}

		// Find next separator and return token
		if i := bytes.Index(data, searchBytes); i >= 0 {
			return i + searchLen, data[0:i], nil
		}

		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return dataLen, data, nil
		}

		// Request more data.
		return 0, nil, nil
	}
}

func ReadLines(path string) (lines []string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return lines, err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}

func ReadRows(path, delimiter string) (rows []string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return rows, err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	scanner.Split(splitOnDelimiter(delimiter))

	for scanner.Scan() {
		rows = append(rows, scanner.Text())
	}

	if len(rows) == 0 {
		err = errors.New("empty file")
	}

	return rows, err
}

func BatchedRows(path, innerDelimiter, outerDelimiter string) (batchedRows [][]string, err error) {
	rows, err := ReadRows(path, outerDelimiter)
	if err != nil {
		return batchedRows, err
	}

	for _, row := range rows {
		batchedRows = append(batchedRows, strings.Split(row, innerDelimiter))
	}
	return batchedRows, nil
}
