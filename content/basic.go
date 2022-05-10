package content

import (
	"encoding/json"
	"time"

	"github.com/micheam/notes"
)

type Basic struct {
	parent    notes.BookID
	id        ContentID
	title     notes.Title
	body      string // TOOD: Introduce Markdown format string types.
	createdAt time.Time
	updatedAt time.Time
}

var _ json.Marshaler = (*Basic)(nil)
var _ json.Unmarshaler = (*Basic)(nil)

func NewBasic(p notes.BookID, t notes.Title) (*Basic, error) {
	if err := t.Validate(); err != nil {
		return nil, err
	}
	id := NewContentID()
	now := time.Now()
	return &Basic{
		parent:    p,
		id:        id,
		title:     t,
		body:      "",
		createdAt: now,
		updatedAt: now,
	}, nil
}

func (c Basic) Parent() notes.BookID { return c.parent }
func (c Basic) ID() ContentID        { return c.id }
func (c Basic) Title() notes.Title   { return c.title }
func (c Basic) Body() string         { return c.body }
func (c Basic) CreatedAt() time.Time { return c.createdAt }
func (c Basic) UpdatedAt() time.Time { return c.updatedAt }

func (c *Basic) SetTitle(t notes.Title) error {
	if err := t.Validate(); err != nil {
		return err
	}
	c.title = t
	c.updatedAt = time.Now()
	return nil
}

func (c *Basic) SetBody(b string) error {
	c.body = b
	c.updatedAt = time.Now()
	return nil
}

func (c Basic) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"id":         c.id.String(),
		"parent":     c.parent.String(),
		"title":      c.title.String(),
		"body":       c.body,
		"created_at": c.createdAt.Format(time.RFC3339),
		"updated_at": c.updatedAt.Format(time.RFC3339),
	}
	return json.Marshal(m)
}

func (c *Basic) UnmarshalJSON([]byte) error {
	panic("unimplemented yet")
}
