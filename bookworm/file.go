package bookworm

import (
	"fmt"
)

type File struct {
	body           string
	totalWordCount int
	wordCounts     map[string]int
}

func Parse() File {
	panic("Not implemented")
}
