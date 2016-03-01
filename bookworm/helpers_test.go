package bookworm

import (
	"testing"
)

func TestCompact(t *testing.T) {
	testCases := []struct {
		input    []string
		expected []string
	}{

		{
			[]string{"", " ", "A", "B"},
			[]string{" ", "A", "B"},
		},
		{
			[]string{"A", " ", ""},
			[]string{"A", " "},
		},
		{
			[]string{"A bility", " ", ""},
			[]string{"A bility", " "},
		},
	}

	for _, testCase := range testCases {
		output := Compact(testCase.input)
		if testCase.expected != output {
			t.Errorf("Unexpected output: Expected: %s, Got: %s", testCase.expected, output)
		}
	}
}
