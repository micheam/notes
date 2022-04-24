package notes

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Page struct {
	id        PageID
	value     string
	createdAT time.Time
	updatedAT time.Time
}

func NewPageID() PageID {
	v := fmt.Sprintf("page:%v", uuid.New())
	return PageID(v)
}

type PageID string

func NewPage(v string) (*Page, error) {
	if len(v) == 0 {
		return nil, errors.New("empty value")
	}
	now := time.Now()
	return &Page{
		id:        NewPageID(),
		value:     v,
		createdAT: now,
		updatedAT: now,
	}, nil
}

func (p Page) String() string {
	text := ellipsis(p.value, 20, "...")
	text = strings.ReplaceAll(text, "\n", " ")
	return fmt.Sprintf("%s (%v)", text, p.id)
}
