package sqlite_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jmoiron/sqlx"
	"github.com/micheam/notes"
	"github.com/micheam/notes/content"
	"github.com/micheam/notes/internal/persistence/sqlite"
)

func initdb(t *testing.T) *sqlx.DB {
	t.Helper()

	db, err := sqlite.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
		return nil
	}
	t.Cleanup(func() { _ = db.Close() })

	// init schema
	err = sqlite.Migrate(db)
	if err != nil {
		t.Fatal("failed to migrate database", err)
		return nil
	}

	return db
}

var opts = cmp.Options{
	cmpopts.IgnoreFields(notes.Book{}, "UpdatedAt"),
	cmpopts.IgnoreFields(content.Basic{}, "UpdatedAt"),
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

func newGeneralContent(t *testing.T, bookID notes.BookID, title notes.Title) *content.Basic {
	cont, err := content.NewBasic(bookID, title)
	if err != nil {
		t.Fatal(err)
	}
	return cont
}

func prepareGeneralCont(t *testing.T, db *sqlx.DB, cont *content.Basic) {
	t.Helper()
	err := sqlite.NewGeneralContentAccess(db).Insert(context.Background(), cont)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.MustExec("delete from general where id=?", cont.ID) })
}
