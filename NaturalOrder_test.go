package main

import (
	"slices"
	"testing"
)

// The reason the natural order exists: the lexicographic order of os.ReadDir
// puts "ep 10" before "ep 2" and shifts every episode number of the folder.
func TestSortNaturallyOrdersNumbersByValue(t *testing.T) {
	fileNames := []string{"ep 10.mkv", "ep 2.mkv", "ep 1.mkv", "ep 20.mkv", "ep 3.mkv"}
	sortNaturally(fileNames)

	expected := []string{"ep 1.mkv", "ep 2.mkv", "ep 3.mkv", "ep 10.mkv", "ep 20.mkv"}
	if !slices.Equal(fileNames, expected) {
		t.Errorf("Expected %v, but got %v", expected, fileNames)
	}
}

// A name can hold several numbers, and the ones before the episode must not
// decide the order on their own.
func TestSortNaturallyComparesEveryNumberOfAName(t *testing.T) {
	fileNames := []string{"Show 2 - ep 10.mkv", "Show 10 - ep 2.mkv", "Show 2 - ep 2.mkv"}
	sortNaturally(fileNames)

	expected := []string{"Show 2 - ep 2.mkv", "Show 2 - ep 10.mkv", "Show 10 - ep 2.mkv"}
	if !slices.Equal(fileNames, expected) {
		t.Errorf("Expected %v, but got %v", expected, fileNames)
	}
}

// A video and its subtitle share everything but the extension: they must stay
// next to each other, whatever the numbers around them.
func TestSortNaturallyKeepsCompanionFilesTogether(t *testing.T) {
	fileNames := []string{"ep 10.srt", "ep 2.mkv", "ep 10.mkv", "ep 2.srt"}
	sortNaturally(fileNames)

	expected := []string{"ep 2.mkv", "ep 2.srt", "ep 10.mkv", "ep 10.srt"}
	if !slices.Equal(fileNames, expected) {
		t.Errorf("Expected %v, but got %v", expected, fileNames)
	}
}

func TestCompareNatural(t *testing.T) {
	cases := []struct {
		a     string
		b     string
		order string
	}{
		{"ep 2.mkv", "ep 10.mkv", "<"},
		{"ep 10.mkv", "ep 2.mkv", ">"},
		{"ep 2.mkv", "ep 2.mkv", "="},
		// A number is only worth what it spells out: the zeros in front of it
		// change nothing, they only settle the tie.
		{"ep 007.mkv", "ep 8.mkv", "<"},
		{"ep 7.mkv", "ep 007.mkv", "<"},
		// Longer than any int64, and still ordered by value.
		{"ep 99999999999999999999.mkv", "ep 100000000000000000000.mkv", "<"},
		// Without digits to compare, the byte order of os.ReadDir is kept.
		{"cover.jpg", "poster.jpg", "<"},
		{"Show.mkv", "show.mkv", "<"},
		{"ep 2.mkv", "ep 2.srt", "<"},
		// A digit run against something else, and a name that stops early.
		{"ep 2.mkv", "ep two.mkv", "<"},
		{"ep 2", "ep 2.mkv", "<"},
	}

	for _, testCase := range cases {
		got := compareNatural(testCase.a, testCase.b)
		if !matchesOrder(got, testCase.order) {
			t.Errorf("Expected %q %s %q, but compareNatural returned %d",
				testCase.a, testCase.order, testCase.b, got)
		}
	}
}

func matchesOrder(got int, order string) bool {
	switch order {
	case "<":
		return got < 0
	case ">":
		return got > 0
	default:
		return got == 0
	}
}
