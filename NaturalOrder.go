package main

import (
	"slices"
	"strings"
)

// sortNaturally orders file names the way a human numbers episodes: the digit
// runs inside a name are compared as numbers, so "ep 2" comes before "ep 10".
// os.ReadDir only offers the lexicographic order, which puts "ep 10" first and
// hands the wrong episode number to every file of a folder whose numbers are
// not zero-padded.
func sortNaturally(fileNames []string) {
	slices.SortFunc(fileNames, compareNatural)
}

// compareNatural compares two names chunk by chunk, a chunk being either a run
// of digits — compared by value — or a single byte, compared as os.ReadDir
// would. It returns a negative number, zero, or a positive number, as the
// convention of slices.SortFunc requires.
func compareNatural(a string, b string) int {
	for len(a) > 0 && len(b) > 0 {
		if isDigit(a[0]) && isDigit(b[0]) {
			aDigits, aRest := splitDigits(a)
			bDigits, bRest := splitDigits(b)
			if order := compareNumbers(aDigits, bDigits); order != 0 {
				return order
			}
			a, b = aRest, bRest

			continue
		}

		if a[0] != b[0] {
			return int(a[0]) - int(b[0])
		}
		a, b = a[1:], b[1:]
	}

	return len(a) - len(b)
}

// compareNumbers compares two runs of digits by the value they spell out. The
// runs can be arbitrarily long, so they are compared as numbers written down
// rather than converted: a file name may hold more digits than an int can.
func compareNumbers(a string, b string) int {
	// "007" and "7" are the same episode. Leading zeros are dropped for the
	// comparison, then used as a tie-breaker so that the order stays total.
	trimmedA := strings.TrimLeft(a, "0")
	trimmedB := strings.TrimLeft(b, "0")

	if len(trimmedA) != len(trimmedB) {
		return len(trimmedA) - len(trimmedB)
	}
	if order := strings.Compare(trimmedA, trimmedB); order != 0 {
		return order
	}

	return len(a) - len(b)
}

// splitDigits cuts the leading run of digits off a name and returns it with
// what is left.
func splitDigits(name string) (digits string, rest string) {
	end := 0
	for end < len(name) && isDigit(name[end]) {
		end++
	}

	return name[:end], name[end:]
}

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}
