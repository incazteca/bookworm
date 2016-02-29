package bookworm

import (
	"os"
	"reflect"
	"testing"
)

type testCase struct {
	expected FileStats
	filter   string
	filePath string
}

func TestParse(t *testing.T) {
	bodyText := "The quick brown fox jumps over the lazy dog"
	wordCounts := map[string]int{
		"the":   2,
		"quick": 1,
		"brown": 1,
		"fox":   1,
		"jumps": 1,
		"over":  1,
		"lazy":  1,
		"dog":   1,
	}

	filteredBodyText := "quick brown fox jumps over lazy dog"
	filteredWordCounts := map[string]int{
		"quick": 1,
		"brown": 1,
		"fox":   1,
		"jumps": 1,
		"over":  1,
		"lazy":  1,
		"dog":   1,
	}

	filteredOBodyText := "The quick jumps the lazy"
	filteredOWordCounts := map[string]int{
		"the":   2,
		"quick": 1,
		"jumps": 1,
		"lazy":  1,
	}

	testCases := []testCase{
		{
			FileStats{bodyText, 9, wordCounts},
			"",
			"../sample_files/fox.txt",
		},
		{
			FileStats{filteredBodyText, 7, filteredWordCounts},
			"the",
			"../sample_files/fox.txt",
		},
		{
			FileStats{filteredBodyText, 7, filteredWordCounts},
			"The",
			"../sample_files/fox.txt",
		},
		{
			FileStats{filteredBodyText, 7, filteredWordCounts},
			"THE",
			"../sample_files/fox.txt",
		},
		{
			FileStats{filteredBodyText, 7, filteredWordCounts},
			"tHe",
			"../sample_files/fox.txt",
		},
		{
			FileStats{filteredBodyText, 7, filteredWordCounts},
			"tHe",
			"../sample_files/fox.txt",
		},
		{
			FileStats{filteredOBodyText, 5, filteredOWordCounts},
			"o",
			"../sample_files/fox.txt",
		},
	}

	for _, testCase := range testCases {
		fh, err := os.Open(testCase.filePath)

		if err != nil {
			t.Fatalf("Unable to open file: %s", testCase.filePath)
		}

		fileStat, err := Parse(fh, testCase.filter)

		if err != nil {
			t.Fatalf("Error in parsing file: %s", err)
		}

		// Compare bodies
		if testCase.expected.Body != fileStat.Body {
			t.Errorf("Expected: %s, Received: %s", testCase.expected.Body, fileStat.Body)
		}

		// Compare total word count
		if testCase.expected.TotalWordCount != fileStat.TotalWordCount {
			t.Errorf("Expected: %d, Received: %d", testCase.expected.TotalWordCount, fileStat.TotalWordCount)
		}

		// Compare word counts
		eq := reflect.DeepEqual(testCase.expected.WordCounts, fileStat.WordCounts)

		if !eq {
			t.Errorf("Expected: %v, Received %v", testCase.expected.WordCounts, fileStat.WordCounts)
		}
	}
}

func TestParseFailNilFile(t *testing.T) {
	_, err := Parse(nil, "")

	if err == nil {
		t.Errorf("Expected error due to nil file handle. But error was not received")
	}
}

func TestfilterWordRegex(t *testing.T) {
	testCases := []struct {
		input    string
		filter   string
		expected string
	}{
		{"Blueberry", "", "Blueberry"},
		{"Blueberry", "Blueberry", ""},
		{"Blueberry", "Blue", ""},
		{"Blueberry", "blue", ""},
		{"Blueberry", "berry", ""},
		{"Blueberry", "berry", ""},
		{"Apple", "berry", "Apple"},
		{"Apple", "ap", ""},
		{"Apple", "", "Apple"},
	}

	for _, testCase := range testCases {
		output, err := filterWord(testCase.input, testCase.filter)

		if err != nil {
			t.Fatalf(err.Error())
		}

		if testCase.expected != output {
			t.Errorf("Results don't match, expected: %s, received: %s", testCase.expected, output)
		}
	}
}

func TestcompactLine(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"This is a line", "This is a line"},
		{"This   is  a   line", "This is a line"},
		{"This is a line  ", "This is a line"},
		{"  This is a line", "This is a line"},
		{"  This is a line  ", "This is a line"},
		{"  This   is   a line  ", "This is a line"},
	}

	for _, testCase := range testCases {
		output := compactLine(testCase.input)

		if testCase.expected != output {
			t.Errorf("Results don't match, expected: %s, received: %s", testCase.expected, output)
		}
	}
}

func TestcleanWord(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"foobar", "foobar"},
		{"foo-bar", "foo-bar"},
		{"\"foobar\"", "foobar"},
		{"O'Malley", "O'Malley"},
		{"tilt-a-whirl", "tilt-a-whirl"},
	}

	for _, testCase := range testCases {
		output, err := cleanWord(testCase.input)

		if err != nil {
			t.Fatalf(err.Error())
		}

		if testCase.expected != output {
			t.Errorf("Results don't match, expected: %s, received: %s", testCase.expected, output)
		}
	}
}
