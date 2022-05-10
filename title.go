package notes

import (
	"errors"
	"fmt"
)

type Title string

var ErrInvalidTitle = errors.New("invalid title")

func (t Title) Validate() error {
	if len(t) == 0 {
		return fmt.Errorf("%w: empty title", ErrInvalidTitle)
	}
	return nil
}

func (t Title) String() string {
	return string(t)
}
