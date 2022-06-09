package sqlite_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/micheam/notes"
	"github.com/micheam/notes/sqlite"
	"github.com/stretchr/testify/assert"
)

func TestContentAccess_Insert(t *testing.T) {
	t.Parallel()
	db := initdb(t)
	ctx := context.Background()

	// prepare book
	book := newBook(t, "some book title")
	prepareBook(t, db, book)

	// exercise
	sut := sqlite.NewContentAccess(db)
	cont := genContent(t, book.ID, "some content title")
	err := sut.Insert(ctx, cont)
	assert.NoError(t, err)
}

func TestContentAccess_Get(t *testing.T) {
	t.Parallel()
	db := initdb(t)
	ctx := context.Background()
	sut := sqlite.NewContentAccess(db)

	// prepare data
	book := newBook(t, "some book title")
	prepareBook(t, db, book)
	cont := genContent(t, book.ID, "some content title")
	prepareContent(t, db, cont)

	got, err := sut.Get(ctx, cont.ID)
	if assert.NoError(t, err) {
		if diff := cmp.Diff(cont, got, opts...); diff != "" {
			t.Errorf("got content mismatch (-want, +got):%s\n", diff)
		}
	}
}

func TestContentAccess_Update(t *testing.T) {
	t.Parallel()
	db := initdb(t)
	ctx := context.Background()

	// Prepare book
	book := newBook(t, "some book title")
	prepareBook(t, db, book)

	// Prepare content
	cont := genContent(t, book.ID, "some content title")
	orgCreatedAt := time.Now().Add(-10 * time.Second)
	cont.CreatedAt = orgCreatedAt
	cont.UpdatedAt = orgCreatedAt
	prepareContent(t, db, cont)

	// Exercise - Edit and Save content
	cont.Title = cont.Title + " updated"
	sut := sqlite.NewContentAccess(db)
	err := sut.Update(ctx, cont)
	assert.NoError(t, err)

	// Verify
	got, err := sut.Get(ctx, cont.ID)
	if assert.NoError(t, err) {
		if diff := cmp.Diff(cont, got, opts...); diff != "" {
			t.Errorf("got cont mismatch (-want, +got):%s\n", diff)
		}
		assert.True(
			t, got.UpdatedAt.After(orgCreatedAt),
			"UpdatedAt: got %v, org %v", got.UpdatedAt, orgCreatedAt,
		)
	}
}

func TestContentAccess_Delete(t *testing.T) {
	t.Parallel()
	db := initdb(t)

	// Prepare data
	book := newBook(t, "test book: "+uuid.New().String())
	prepareBook(t, db, book)
	cont := genContent(t, book.ID, "some content title")
	prepareContent(t, db, cont)

	// Exercise
	sut := sqlite.NewContentAccess(db)
	t.Run("must delete content", func(t *testing.T) {
		err := sut.Delete(context.Background(), cont)
		if assert.NoError(t, err) {
			err = db.Get(&struct{}{}, `select TRUE from general where id=?`, cont.ID)
			assert.ErrorIs(t, err, sql.ErrNoRows)
		}
	})
	t.Run("no error if content not exists", func(t *testing.T) {
		title := notes.Title("no-exists: " + uuid.New().String())
		cont := genContent(t, book.ID, title)
		err := sut.Delete(context.Background(), cont)
		assert.NoError(t, err)
	})
}

func TestContentAccess_List(t *testing.T) {
	t.Parallel()
	db := initdb(t)

	// Prepare data
	book := newBook(t, "test book: "+uuid.New().String())
	prepareBook(t, db, book)
	for i := 0; i < 100; i++ {
		title := notes.Title("title" + uuid.New().String())
		prepareContent(t, db, genContent(t, book.ID, title))
	}

	// Exercise
	sut := sqlite.NewContentAccess(db)
	list, err := sut.List(context.Background())
	if assert.NoError(t, err) {
		assert.Len(t, list, 100)
	}
}
