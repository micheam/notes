package notes

import (
	"time"
)

type GeneralContent struct {
	Parent    BookID
	ID        ContentID
	Title     Title
	Body      string // TOOD: Introduce Markdown format string types.
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewGeneralContent(p BookID, t Title) (*GeneralContent, error) {
	if err := t.Validate(); err != nil {
		return nil, err
	}
	id := NewContentID()
	now := time.Now()
	return &GeneralContent{
		Parent:    p,
		ID:        id,
		Title:     t,
		Body:      "",
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
