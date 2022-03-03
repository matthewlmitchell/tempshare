package mysql

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func newTestDatabase(t *testing.T) (*sql.DB, func()) {
	db, err := sql.Open("mysql", fmt.Sprintf("test_web:%s@/test_tempshare?parseTime=true&multiStatements=true", os.Getenv("TEMPSHARE_SQL_AUTH_TEST")))
	if err != nil {
		t.Fatal(err)
	}

	setupScript, err := ioutil.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(setupScript))
	if err != nil {
		t.Fatal(err)
	}

	return db, func() {
		teardownScript, err := ioutil.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}

		_, err = db.Exec(string(teardownScript))
		if err != nil {
			t.Fatal(err)
		}

		db.Close()
	}
}
