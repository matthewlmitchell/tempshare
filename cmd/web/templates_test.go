package main

import (
	"log"
	"testing"
	"time"
)

func TestFormattedDate(t *testing.T) {

	testCases := []struct {
		testName       string
		testInput      time.Time
		expectedOutput string
	}{
		{
			testName:       "Empty time",
			testInput:      time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
			expectedOutput: "",
		},
		{
			testName:       "Valid time",
			testInput:      time.Date(2022, 2, 21, 12, 30, 0, 0, time.UTC),
			expectedOutput: "Feb 21 2022 at 12:30",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			realOutput := FormattedDate(testCase.testInput)

			if realOutput != testCase.expectedOutput {
				log.Printf("Expected %s, received %s\n", testCase.expectedOutput, realOutput)
			}
		})
	}
}
