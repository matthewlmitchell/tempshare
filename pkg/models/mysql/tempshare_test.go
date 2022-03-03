package mysql

import (
	"crypto/sha256"
	"reflect"
	"testing"
	"time"

	"github.com/matthewlmitchell/tempshare/pkg/models"
)

func TestNew(t *testing.T) {

	type tempShareInput struct {
		Text      string
		Expires   string
		ViewLimit string
	}

	testCases := []struct {
		name              string
		inputTempShare    tempShareInput
		expectedTempShare *models.TempShare
		expectedError     error
	}{
		{
			name: "Valid input",
			inputTempShare: tempShareInput{
				Text:      "This is an example tempshare for testing purposes!",
				Expires:   "7",
				ViewLimit: "10",
			},
			expectedTempShare: &models.TempShare{
				Text:      "This is an example tempshare for testing purposes!",
				ViewLimit: 10,
			},
			expectedError: nil,
		},
		{
			name: "Empty text",
			inputTempShare: tempShareInput{
				Text:      "",
				Expires:   "1",
				ViewLimit: "1",
			},
			expectedTempShare: &models.TempShare{
				Text:      "",
				ViewLimit: 1,
			},
			expectedError: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			db, teardown := newTestDatabase(t)
			defer teardown()

			model := &TempShareModel{db}

			tempShare, err := model.New(testCase.inputTempShare.Text, testCase.inputTempShare.Expires, testCase.inputTempShare.ViewLimit)
			if err != testCase.expectedError {
				t.Errorf("Expected %v, received %v", testCase.expectedError, err)
			}

			if tempShare.ViewLimit != testCase.expectedTempShare.ViewLimit {
				t.Errorf("Expected %d, received %d", testCase.expectedTempShare.ViewLimit, tempShare.ViewLimit)
			}

		})
	}
}

func TestGet(t *testing.T) {

	testCases := []struct {
		name                string
		inputPlainTextToken string
		expectedTempShare   *models.TempShare
		expectedError       error
	}{
		{
			name:                "Valid Get",
			inputPlainTextToken: "FTR43TPBEWDCQ4B2HRCNXPSDBXFEAQ44QWC7QZ2P5D5NW3Y64UJA",
			expectedTempShare: &models.TempShare{
				Text:      "This is an example tempshare for testing purposes!",
				Created:   time.Date(2022, 3, 2, 12, 0, 0, 0, time.UTC),
				Expires:   time.Date(2048, 3, 9, 12, 0, 0, 0, time.UTC),
				Views:     0,
				ViewLimit: 1,
			},
			expectedError: nil,
		},
		{
			name:                "No matching record",
			inputPlainTextToken: "FEAQ44QWC7QZ2P5D5NW3Y64UJFTR43TPBEWDCQ4B2HRCNXPSDBXA",
			expectedTempShare:   nil,
			expectedError:       models.ErrNoRecord,
		},
		{
			name:                "Empty token",
			inputPlainTextToken: "",
			expectedTempShare:   nil,
			expectedError:       models.ErrNoRecord,
		},
		{
			name:                "Expired tempshare",
			inputPlainTextToken: "HVN2JMTD5DVPODS632YXWVT6REYSXR26O7B3G5ZBQRD72IOBYTVA",
			expectedTempShare:   nil,
			expectedError:       models.ErrNoRecord,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			if testCase.expectedTempShare != nil {
				hash := sha256.Sum256([]byte(testCase.inputPlainTextToken))
				testCase.expectedTempShare.URLToken = hash[:]
			}

			db, teardown := newTestDatabase(t)
			defer teardown()

			model := &TempShareModel{db}

			tempShare, err := model.Get(testCase.inputPlainTextToken)
			if err != testCase.expectedError {
				t.Errorf("Expected %v, received %v", testCase.expectedError, err)
			}

			if !reflect.DeepEqual(tempShare, testCase.expectedTempShare) {
				t.Errorf("Expected %v, received %v", testCase.expectedTempShare, tempShare)
			}

			// Testing logic ...

		})
	}
}
