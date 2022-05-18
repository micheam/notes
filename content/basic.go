package content

import (
	"time"

	"github.com/micheam/notes"
)

type Basic struct {
	Parent    notes.BookID
	ID        ContentID
	Title     notes.Title
	Body      string // TOOD: Introduce Markdown format string types.
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewBasic(p notes.BookID, t notes.Title) (*Basic, error) {
	if err := t.Validate(); err != nil {
		return nil, err
	}
	id := NewContentID()
	now := time.Now()
	return &Basic{
		Parent:    p,
		ID:        id,
		Title:     t,
		Body:      "",
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
