package testing

import (
	"testing"
	"time"

	"github.com/micheam/notes"
)

func RandBook(t *testing.T) *notes.Book {
	t.Helper()
	return &notes.Book{
		ID:        notes.NewBookID(),
		Title:     RandStr(10),
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
}

func RandContent(t *testing.T) *notes.Content {
	t.Helper()
	return &notes.Content{
		ID:        notes.NewContentID(),
		Title:     RandTitle(t),
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}
}

func RandTitle(t *testing.T) notes.Title {
	t.Helper()
	return notes.Title("title: " + RandStr(20))
}
