package sqlite_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jmoiron/sqlx"
	"github.com/micheam/notes"
	"github.com/micheam/notes/internal/persistence/sqlite"
)

func initdb(t *testing.T) *sqlx.DB {
	t.Helper()

	dir := t.TempDir()
	db, err := sqlx.Open("sqlite3", filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatal(err)
		return nil
	}
	t.Cleanup(func() { _ = db.Close() })

	// init schema
	schema, err := os.ReadFile("schema.sql")
	if err != nil {
		t.Fatal(err)
		return nil
	}
	_ = db.MustExec(string(schema))

	return db
}

var opts = cmp.Options{
	cmpopts.IgnoreFields(notes.Book{}, "UpdatedAt"),
	cmpopts.IgnoreUnexported(notes.Book{}),
	cmpopts.EquateApproxTime(1 * time.Second),
}

func newBook(t *testing.T, title string) *notes.Book {
	book, err := notes.NewBook(title)
	if err != nil {
		t.Fatal(err)
	}
	return book
}

func prepareBook(t *testing.T, db *sqlx.DB, book *notes.Book) {
	t.Helper()
	err := sqlite.NewBookAccess(db).InsertBook(context.Background(), book)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.MustExec("delete from book where id=?", book.ID) })
}
