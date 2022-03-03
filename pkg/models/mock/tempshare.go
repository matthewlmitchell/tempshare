package mock

import (
	"crypto/sha256"
	"time"

	"github.com/matthewlmitchell/tempshare/pkg/models"
)

// This is a mock version of the TempShareModel{DB: *sql.DB} struct
type TempShareModel struct{}

var mockTempShare = &models.TempShare{
	Text:      "This is an example tempshare for testing purposes!",
	PlainText: "MUPPH5PDKV7AGCUAAEERL5ARIXICVVGYLRIV365X5XSV3EKISAXQ",
	Created:   time.Now(),
	Expires:   time.Now().Add(24 * time.Hour),
	Views:     0,
	ViewLimit: 1,
}

func (model *TempShareModel) New(text string, expires string, viewlimit string) (*models.TempShare, error) {

	hash := sha256.Sum256([]byte(mockTempShare.PlainText))
	mockTempShare.URLToken = hash[:]

	// TODO: Insert(...)

	return mockTempShare, nil
}

func (model *TempShareModel) Insert(urlToken []byte, text string, expires string, viewlimit int) error {

	return nil
}
func (model *TempShareModel) Get(plaintextToken string) (*models.TempShare, error) {

	if plaintextToken == "MUPPH5PDKV7AGCUAAEERL5ARIXICVVGYLRIV365X5XSV3EKISAXQ" {
		// TODO: Update(...)

		return mockTempShare, nil
	}

	return nil, models.ErrNoRecord
}

func (model *TempShareModel) Update(plaintextToken string) error {

	if plaintextToken == "MUPPH5PDKV7AGCUAAEERL5ARIXICVVGYLRIV365X5XSV3EKISAXQ" {
		return nil
	}

	return models.ErrNoRecord
}
