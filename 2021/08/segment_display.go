package main

import (
	"flag"
	"math"
	"os"
	"strings"

	"github.com/aaronireland/advent-of-code/pkg/input"
	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
)

var (
	filePath         string
	segmentLengthMap = map[int][]int{
		2: {1},
		3: {7},
		4: {4},
		5: {2, 3, 5},
		6: {0, 6, 9},
		7: {8},
	}
)

func init() {
	log.SetHandler(text.New(os.Stderr))
	flag.StringVar(&filePath, "input", "./input.txt", "The input data text file path")
}

// filter strips the filter string from the input string if and only each character in the filter occurs in the input,
// regardless of order
func filter(input, filter string) string {
	filtered := new(string)
	*filtered = input
	for _, char := range filter {
		*filtered = strings.ReplaceAll(*filtered, string(char), "")
	}
	if len(*filtered)+len(filter) == len(input) {
		return *filtered
	}
	return input
}

// scrambled I/O signals aren't guaranteed to be an any certain order so a normal key lookup won't work
func getMappedDigit(digitMap map[string]int, segments string) (int, bool) {
	for seg, num := range digitMap {
		if len(filter(segments, seg)) == 0 && len(filter(seg, segments)) == 0 {
			return num, true
		}
	}

	return 0, false
}

// makeSegmentMap decrypts the scrambled I/O signals. Decimal digits 1, 4, 7, and 8 are generated using a unique number
// of display segments. For the less trivial case of distinguishing {2, 3, 5} and {0, 6, 9} which have 5 and 6 displays
// segments respectively, we need a few extra filters: namely the right two segments that make up a 1 (called c & f), &
// the extra pair of segments (b & d) that make up a 4.
func makeSegmentMap(display []string, lengthOnly bool) map[string]int {
	var (
		segmentsToDigit = make(map[string]int)
		cf, bcdf, bd    string // decimal 1, decimal 4, and filtering 1 from 4 identifies the bd segment pair
	)
	for _, ioSegments := range display {
		for _, digit := range strings.Fields(ioSegments) {
			if mappedDigit, ok := segmentLengthMap[len(digit)]; ok {
				if len(mappedDigit) == 1 {
					if _, ok := getMappedDigit(segmentsToDigit, digit); !ok {
						segmentsToDigit[digit] = mappedDigit[0]
						switch mappedDigit[0] {
						case 1:
							cf = digit
						case 4:
							bcdf = digit
						}
					}
				}
			}
		}
	}

	if lengthOnly {
		return segmentsToDigit
	}

	bd = filter(bcdf, cf) // map the b and d segments by removing the mapped "cd" segments from bcdf (digit 4)

	for _, ioSegments := range display {
		for _, digit := range strings.Fields(ioSegments) {
			if len(digit) == 5 {
				if len(digit) == len(filter(digit, cf))+len(cf) { // only decimal 3 has 5 segments and contains the cf segment pair
					segmentsToDigit[digit] = 3
				} else {
					if len(digit) == len(filter(digit, bd))+len(bd) { // only decimal 5 has 5 segments and contains the bd segment pair
						segmentsToDigit[digit] = 5
					} else {
						segmentsToDigit[digit] = 2 // only decimal 2 have a length of 5 and does not contain either bd, nor cf segment pairs
					}
				}
			}
			if len(digit) == 6 {
				if len(digit) == len(filter(digit, cf)) {
					segmentsToDigit[digit] = 6 // only decimal 6 has 6 segments and does not contain either bd nor cf segment pairs
				} else {
					if len(filter(filter(digit, bd), cf)) == 2 { // only decimal 9 has 6 segments and contains both bd and cf segments pairs
						segmentsToDigit[digit] = 9
					} else {
						segmentsToDigit[digit] = 0 // only decimal 0 has 6 segments and only contains the cf segment pair
					}
				}
			}
		}

	}

	return segmentsToDigit
}

func countDigitsInDisplays(displays [][]string, lengthOnly bool) map[int]int {
	digitCount := make(map[int]int)

	for _, display := range displays {
		digitFor := makeSegmentMap(display, lengthOnly)
		output := display[1]
		for _, displayedSegments := range strings.Fields(output) {
			if digit, ok := getMappedDigit(digitFor, displayedSegments); ok {
				digitCount[digit]++
			}
		}
	}

	return digitCount
}

func sumOutputInDisplays(displays [][]string) (outputSum int64) {
	sumOutputCh := make(chan int64)

	for _, display := range displays {
		go func(d []string, resCh chan int64) {
			var displayed int64
			digitFor := makeSegmentMap(d, false)
			displayedSegments := strings.Fields(d[1])

			for x, segment := range displayedSegments {
				if digit, ok := getMappedDigit(digitFor, segment); ok {
					multiplier := int(math.Pow(10, float64(len(displayedSegments)-(1+x))))
					displayed += int64(digit * multiplier)
				} else {
					log.WithField("segment-map", digitFor).Warnf("Numeric digit not found for %s", segment)
				}
			}
			resCh <- displayed
		}(display, sumOutputCh)
	}

	for i := 0; i < len(displays); i++ {
		outputSum += <-sumOutputCh
	}

	return outputSum
}

func main() {
	flag.Parse()
	data, err := input.BatchedRows(filePath, "|", "\n")
	if err != nil {
		log.WithError(err).Fatalf("Failed to read valid input data from %s", filePath)
	}

	totalCountPartOne := 0 // Count 1, 4, 7, and 8
	for digit, count := range countDigitsInDisplays(data, true) {
		if digit == 1 || digit == 4 || digit == 7 || digit == 8 {
			totalCountPartOne += count
		}
		log.Infof("Digit %d -> %d occurrences", digit, count)
	}
	log.Infof("Total count for digits 1, 4, 7, and 8 is %d", totalCountPartOne)
	totalSumOfOutput := sumOutputInDisplays(data)

	log.Infof("Total sum of all output is %d", totalSumOfOutput)
}
