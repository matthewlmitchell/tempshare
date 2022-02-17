package mysql

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"strconv"
	"time"

	"github.com/matthewlmitchell/tempshare/pkg/models"
)

type TempShareModel struct {
	DB *sql.DB
}

func (model *TempShareModel) New(text string, expires string, viewlimit string) (*models.TempShare, error) {
	maxViews, err := strconv.Atoi(viewlimit)
	if err != nil {
		return nil, err
	}

	expiry, err := strconv.Atoi(expires)
	if err != nil {
		return nil, err
	}

	tempShare, err := generateTempShare(text, expiry, maxViews)
	if err != nil {
		return nil, err
	}

	err = model.Insert(tempShare.URLToken, tempShare.Text, expires, tempShare.ViewLimit)
	return tempShare, err
}

// Insert a string of text into the database with a given expiry
// and return a token string for formatting into a shareable URL
func (model *TempShareModel) Insert(urlToken []byte, text string, expires string, viewlimit int) error {

	sqlStatement := `INSERT INTO texts (urltoken, text, created, expires, views, viewlimit) 
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY), ?, ?)`

	sqlArgs := []interface{}{urlToken, text, expires, 0, viewlimit}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := model.DB.ExecContext(ctx, sqlStatement, sqlArgs...)

	return err
}

func (model *TempShareModel) Get(plaintextToken string) (*models.TempShare, error) {
	urlToken := sha256.Sum256([]byte(plaintextToken))

	sqlStatement := `SELECT urltoken, text, created, expires, views, viewlimit FROM texts
	WHERE expires > UTC_TIMESTAMP() AND views < viewlimit AND urltoken = ?`

	sqlRow := model.DB.QueryRow(sqlStatement, urlToken)

	tempShare := &models.TempShare{}

	err := sqlRow.Scan(&tempShare.URLToken, &tempShare.Text, &tempShare.Created, &tempShare.Expires, &tempShare.Views, &tempShare.ViewLimit)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return tempShare, nil
}

func (model *TempShareModel) Delete(tempShare *models.TempShare) error {

	sqlStatement := `DELETE from texts WHERE urltoken = ?`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := model.DB.ExecContext(ctx, sqlStatement, tempShare.URLToken)

	return err
}

func generateTempShare(text string, expires int, viewlimit int) (*models.TempShare, error) {
	tempShare := &models.TempShare{
		Text:      text,
		Expires:   time.Now().Add(time.Duration(expires*24) * time.Hour),
		ViewLimit: viewlimit,
	}

	randBytes := make([]byte, 32)
	_, err := rand.Read(randBytes)
	if err != nil {
		return nil, err
	}

	tempShare.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randBytes)
	hash := sha256.Sum256([]byte(tempShare.PlainText))
	tempShare.URLToken = hash[:]

	return tempShare, nil
}
