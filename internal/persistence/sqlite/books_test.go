package sqlite_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/micheam/notes"
	"github.com/micheam/notes/internal/persistence/sqlite"
	"github.com/stretchr/testify/assert"

	_ "github.com/mattn/go-sqlite3"
)

// TODO(micheam): want testing helper: RandomBook()
// maybe notes/fazzing package?

func TestBookAccess_InsertBook(t *testing.T) {
	t.Parallel()
	db := initdb(t)
	sut := sqlite.NewBookAccess(db)
	ctx := context.Background()
	book := newBook(t, "some title")
	err := sut.InsertBook(ctx, book)
	assert.NoError(t, err)
}

func TestBookAccess_GetBook(t *testing.T) {
	t.Parallel()
	db := initdb(t)
	book := newBook(t, "some title: TestBookAccess_GetBook")
	prepareBook(t, db, book)

	sut := sqlite.NewBookAccess(db)
	got, err := sut.GetBook(context.Background(), book.ID)
	if assert.NoError(t, err) {
		if diff := cmp.Diff(book, got, opts...); diff != "" {
			t.Errorf("got book mismatch (-want, +got):%s\n", diff)
		}
	}
}

func TestBookAccess_UpdateBook(t *testing.T) {
	t.Parallel()
	db := initdb(t)
	ctx := context.Background()
	book := newBook(t, "some title:")
	orgCreatedAt := time.Now().Add(-10 * time.Second)
	book.CreatedAt = orgCreatedAt
	book.UpdatedAt = orgCreatedAt
	prepareBook(t, db, book)

	// Exercise
	book.Title = book.Title + " updated"
	sut := sqlite.NewBookAccess(db)
	err := sut.UpdateBook(ctx, book)
	assert.NoError(t, err)

	// Verify
	got, err := sut.GetBook(ctx, book.ID)
	if assert.NoError(t, err) {
		if diff := cmp.Diff(book, got, opts...); diff != "" {
			t.Errorf("got book mismatch (-want, +got):%s\n", diff)
		}
		assert.True(
			t, got.UpdatedAt.After(orgCreatedAt),
			"UpdatedAt: got %v, org %v", got.UpdatedAt, orgCreatedAt,
		)
	}
}

func TestBookAccess_ListBooks(t *testing.T) {
	t.Parallel()
	db := initdb(t)

	prepare := func(t *testing.T, title string) *notes.Book {
		book := newBook(t, title)
		prepareBook(t, db, book)
		return book
	}
	books := []*notes.Book{}
	for i := 0; i < 100; i++ {
		title := "test book " + uuid.New().String()
		books = append(books, prepare(t, title))
	}

	sut := sqlite.NewBookAccess(db)
	got, err := sut.ListBooks(context.Background())
	if assert.NoError(t, err) {
		if diff := cmp.Diff(books, got, opts...); diff != "" {
			t.Errorf("list book result mismatch (-want, +got):%s\n", diff)
		}
	}
}

func TestBookAccess_DeleteBook(t *testing.T) {
	t.Parallel()
	db := initdb(t)
	book := newBook(t, "test book: "+uuid.New().String())
	prepareBook(t, db, book)
	sut := sqlite.NewBookAccess(db)
	t.Run("must delete book", func(t *testing.T) {
		err := sut.DeleteBook(context.Background(), book)
		if assert.NoError(t, err) {
			err = db.Get(&struct{}{}, `select TRUE from book where id=?`, book.ID)
			assert.ErrorIs(t, err, sql.ErrNoRows)
		}
	})
	t.Run("no error if book not exists", func(t *testing.T) {
		book := newBook(t, "no-exists: "+uuid.New().String())
		err := sut.DeleteBook(context.Background(), book)
		assert.NoError(t, err)
	})
}
