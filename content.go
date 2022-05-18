package notes

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type ContentID string

func NewContentID() ContentID {
	v := "content:" + uuid.New().String()
	return ContentID(v)
}

var ErrInvalidContentID = errors.New("invalid ContentID")

func ParseContentID(s string) (ContentID, error) {
	pref, rowuuid, found := strings.Cut(s, ":")
	if !found {
		return "", fmt.Errorf("%w: missing prefix", ErrInvalidContentID)
	}
	if pref != "content" {
		return "", fmt.Errorf("%w: wrong prefix", ErrInvalidContentID)
	}
	_, err := uuid.Parse(rowuuid)
	if err != nil {
		return "", fmt.Errorf("parse uuid: %w", err)
	}
	return ContentID("content:" + rowuuid), nil
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
