package main

import (
	"testing"
)

func TestLineCounter(t *testing.T) {
	testCases := []struct {
		name     string
		filename string
		expected int
	}{
		{name: "Positive test case", filename: "test/test.txt", expected: 0},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			file, err := openFile(tc.filename)
			if err != nil {
				t.Fatal(err)
			}
			got, err := LineCounter(file)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.expected {
				t.Errorf("Expected %q got %q instead", tc.expected, got)
			}

		})

	}
}
