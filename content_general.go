package notes

import (
	"time"
)

type Basic struct {
	Parent    BookID
	ID        ContentID
	Title     Title
	Body      string // TOOD: Introduce Markdown format string types.
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewBasic(p BookID, t Title) (*Basic, error) {
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
