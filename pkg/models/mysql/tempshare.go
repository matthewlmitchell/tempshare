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

// New generates a sha256 hash of a supplied text string, then calls Insert to create a
// new entry in our SQL database. After insertion, a *models.TempShare struct is returned
// containing the necessary info for retrieving the data from the SQL db.
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
// The primary key is a sha256 hash used as a token in the URL
// e.g. /view?token=XXXXXX
func (model *TempShareModel) Insert(urlToken []byte, text string, expires string, viewlimit int) error {

	sqlStatement := `INSERT INTO texts (urltoken, text, created, expires, views, viewlimit) 
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY), ?, ?)`

	sqlArgs := []interface{}{urlToken, text, expires, 0, viewlimit}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := model.DB.ExecContext(ctx, sqlStatement, sqlArgs...)

	return err
}

// Get accepts a base32 encoded string as a primary key and retrieves the corresponding
// entry from our SQL database if it exists (and if it is not expired/exceeding view limits).
// The data is scanned into a models.TempShare{} struct and returned,
// the view count of the DB entry is then incremented to reflect that the data has been accessed
func (model *TempShareModel) Get(plaintextToken string) (*models.TempShare, error) {

	sqlStatement := `SELECT urltoken, text, created, expires, views, viewlimit FROM texts
	WHERE expires > UTC_TIMESTAMP() AND views < viewlimit AND urltoken = ?`

	urlTokenHash := sha256.Sum256([]byte(plaintextToken))
	urlToken := urlTokenHash[:]

	sqlRow := model.DB.QueryRow(sqlStatement, urlToken)

	tempShare := &models.TempShare{}

	err := sqlRow.Scan(&tempShare.URLToken, &tempShare.Text, &tempShare.Created, &tempShare.Expires, &tempShare.Views, &tempShare.ViewLimit)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	// Increment the number of views for the MySQL record
	err = model.Update(plaintextToken)
	if err == models.ErrNoRecord {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return tempShare, nil
}

// Update accepts a string (which should be base32 encoded), which is our primary key
// after taking a sha256 hash, and attempts to increment the view count of the
// corresponding row in our SQL database.
func (model *TempShareModel) Update(plaintextToken string) error {

	sqlStatement := `UPDATE texts
	SET views = views + 1 WHERE urltoken = ?`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	urlTokenHash := sha256.Sum256([]byte(plaintextToken))
	urlToken := urlTokenHash[:]

	result, err := model.DB.ExecContext(ctx, sqlStatement, urlToken)
	if err != nil {
		return err
	}

	// If no rows were affected/updated by the sql statement, then the plaintextToken
	// did not correspond to a valid row in our sql database. Return that
	// there was no record matching the request.
	numRowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if numRowsAffected == 0 {
		return models.ErrNoRecord
	}

	return nil
}

// Delete removes the row from our SQL database which corresponds
// to the provided models.TempShare{} struct.
func (model *TempShareModel) Delete(tempShare *models.TempShare) error {

	sqlStatement := `DELETE from texts WHERE urltoken = ?`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := model.DB.ExecContext(ctx, sqlStatement, tempShare.URLToken)
	if err != nil {
		return err
	}

	// If no rows were affected in our table, then there was no record
	// in our SQL database that matches the provided models.TempShare struct.
	numRowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if numRowsAffected == 0 {
		return models.ErrNoRecord
	}

	return nil
}

// generateTempShare accepts a string of text, the number of days before expiry,
// and a maximum view count. These values are parsed into a models.TempShare{}
// struct, a base32 encoded string is randomly generated to be used as a shareable URL,
// and a sha256 hash of the URL token is generated.
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
