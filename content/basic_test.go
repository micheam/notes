package content

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/micheam/notes"
	"github.com/stretchr/testify/assert"
)

func TestBasic_MarshalJSON(t *testing.T) {
	bookID := notes.MustParseBookID("book:573c14a8-d397-4a67-9dd6-d9784ca14b86")
	id := MustParseContentID("content:a9e57e7d-73b9-4ad7-ab27-d291d4bf0de9")
	ts := time.Date(1984, time.February, 9, 13, 0, 45, 0, notes.JST())
	content := &Basic{
		parent:    bookID,
		title:     "this is a title",
		id:        id,
		createdAt: ts,
		updatedAt: ts,
	}
	want := `{` +
		`"body":"",` +
		`"created_at":"1984-02-09T13:00:45+09:00",` +
		`"id":"content:a9e57e7d-73b9-4ad7-ab27-d291d4bf0de9",` +
		`"parent":"book:573c14a8-d397-4a67-9dd6-d9784ca14b86",` +
		`"title":"this is a title",` +
		`"updated_at":"1984-02-09T13:00:45+09:00"` +
		`}`
	got, err := content.MarshalJSON()
	if assert.NoError(t, err) {
		if diff := cmp.Diff([]byte(want), got); diff != "" {
			t.Errorf("json mismatch (-want, +got):%s\n", diff)
		}
	}
}
