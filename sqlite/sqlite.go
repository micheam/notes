package sqlite

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const timefmt = "2006-01-02 15:04:05.999999999Z07:00"

// parseDatetime parse s into time.Time then return result.
//
// May panic, if layout not match with "2006-01-02 15:04:05.999999999Z07:00".
//
// TODO(micheam): Adjust the layout of date/time items to RFC3339.
// This needs to be changed to match the behavior at the time of
// various records.
func parseDatetime(s string) time.Time {
	t, err := time.Parse(timefmt, s)
	if err != nil {
		panic(fmt.Errorf("try to parse: %s: %w", s, err))
	}
	return t
}

// formatDatetime convert t time.Time into string.
func formatDatetime(t time.Time) string {
	return t.Format(timefmt)
}

func Open(path string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return db, nil
}
