package notes

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrContentNotFound = errors.New("content not found")
	ErrBookNotFound    = errors.New("book not found")

	ErrInvalidTitle     = errors.New("invalid title")
	ErrInvalidContentID = errors.New("invalid ContentID")
	ErrInvalidBookID    = errors.New("invalid BookID")
	ErrInvalidArgument  = errors.New("invalid Argument")
)

type ValidationError struct {
	fieldErrors map[string][]error
}

func NewValidationError(field string, err error) *ValidationError {
	return &ValidationError{
		fieldErrors: map[string][]error{field: {err}},
	}
}

var _ error = (*ValidationError)(nil)

func (verr ValidationError) Error() string {
	var sb strings.Builder
	for field, errs := range verr.fieldErrors {
		sb.WriteString(fmt.Sprintf("(%s) ", field))
		sb.WriteString(errorsJoin(errs, ", ").Error())
	}
	return sb.String()
}

func errorsJoin(errs []error, sep string) error {
	switch len(errs) {
	case 0:
		return nil
	case 1:
		return errs[0]
	default:
		var sb strings.Builder
		sb.WriteString(errs[0].Error())
		for i := 1; i < len(errs); i++ {
			sb.WriteString(sep)
			sb.WriteString(errs[i].Error())
		}
		return errors.New(sb.String())
	}
}
