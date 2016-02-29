package bookworm

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strings"
)

type FileStats struct {
	Body           string
	TotalWordCount int
	WordCounts     map[string]int
}

func Parse(src io.Reader, filter string) (FileStats, error) {
	var fileStat FileStats

	if src == nil {
		return fileStat, errors.New("Input file is nil")
	}

	fileStat.WordCounts = make(map[string]int)

	scanner := bufio.NewScanner(src)

	for scanner.Scan() {
		var line string

		if filter != "" {
			filteredLine, err := filterWord(scanner.Text(), filter)

			if err != nil {
				return fileStat, err
			}

			line = compactLine(filteredLine)
		} else {
			line = scanner.Text()
		}

		words := strings.Fields(strings.ToLower(line))
		fileStat.TotalWordCount += len(words)

		// Default value for non-existent keys is 0. So we can just add them up
		for _, word := range words {
			cWord, err := cleanWord(word)

			if err != nil {
				return fileStat, err
			}

			fileStat.WordCounts[cWord] += 1
		}

		if fileStat.Body == "" {
			fileStat.Body = line
		} else {
			fileStat.Body = strings.Join([]string{fileStat.Body, line}, "\n")
		}
	}

	return fileStat, nil
}

func filterWord(src string, filter string) (string, error) {
	if filter == "" {
		return src, nil
	}

	pattern := strings.Join([]string{"((?i)\\w*", filter, "\\w*)"}, "")

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}

	return regex.ReplaceAllString(src, ""), nil
}

// Remove any extra spaces from the line that might have been caused by word removal
func compactLine(line string) string {
	splitStrings := strings.Split(strings.TrimSpace(line), " ")
	return strings.Join(Compact(splitStrings), " ")
}

// Remove any punctuation marks from a word, except for apostrophes and dashes
func cleanWord(word string) (string, error) {
	matchWordRegex, err := regexp.Compile("[\\w'-]+")
	if err != nil {
		return "", err
	}

	cleanWord := matchWordRegex.FindString(word)

	return cleanWord, nil
}
