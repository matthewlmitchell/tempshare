package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord = errors.New("models: no record found matching your request")
)

type TempShare struct {
	Text      string
	PlainText string
	URLToken  []byte
	Created   time.Time
	Expires   time.Time
	Views     int
	ViewLimit int
}
