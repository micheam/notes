package notes

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// =====================================
//
// BookID
//
// =====================================

type BookID string

const prefixBookID = "book:"

func (b BookID) String() string {
	return string(b)
}

func (b BookID) Empty() bool {
	return len(string(b)) == len(prefixBookID)
}

func NewBookID() BookID {
	return BookID(uuid.New().String())
}

func ParseBookID(s string) (BookID, error) {
	pref, rowid, found := strings.Cut(s, ":")
	if !found {
		rowid = s
	} else if pref != prefixBookID {
		return "", fmt.Errorf("%w: wrong prefix", ErrInvalidBookID)
	}
	_, err := uuid.Parse(rowid)
	if err != nil {
		return "", fmt.Errorf("parse uuid: %w", err)
	}
	return BookID(rowid), nil
}

func MustParseBookID(s string) BookID {
	bid, err := ParseBookID(s)
	if err != nil {
		panic(err)
	}
	return bid
}

// =====================================
//
// ContentID
//
// =====================================

type ContentID string

func NewContentID() ContentID {
	v := uuid.New().String()
	return ContentID(v)
}

func ParseContentID(s string) (ContentID, error) {

	pref, rowid, found := strings.Cut(s, ":")
	if found && pref != "content" {
		return "", fmt.Errorf("%w: wrong prefix", ErrInvalidContentID)
	}
	if !found {
		rowid = s
	}
	_, err := uuid.Parse(rowid)
	if err != nil {
		return "", fmt.Errorf("parse uuid: %w: %s", err, rowid)
	}
	return ContentID(rowid), nil
}

func MustParseContentID(s string) ContentID {
	cid, err := ParseContentID(s)
	if err != nil {
		panic(err)
	}
	return cid
}

func (c ContentID) String() string {
	return string(c)
}

// =====================================
//
// Title
//
// =====================================

type Title string

func (t Title) Validate() error {
	if len(t) == 0 {
		return fmt.Errorf("%w: empty title", ErrInvalidTitle)
	}
	return nil
}

func (t Title) String() string {
	return string(t)
}

// =====================================
//
// TimeStamp
//
// =====================================

type TimeStamp struct {
	value time.Time
}

var (
	jstLocation *time.Location
	jstOnce     sync.Once
)

func JST() *time.Location {
	if jstLocation == nil {
		jstOnce.Do(func() {
			l, err := time.LoadLocation("Asia/Tokyo")
			if err != nil {
				l = time.FixedZone("JST2", +9*60*60)
			}
			jstLocation = l
		})
	}
	return jstLocation
}

func (t TimeStamp) Format() string {
	return t.value.Format(time.RFC3339)
}
