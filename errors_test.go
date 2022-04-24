package notes

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestValidationError_Error(t *testing.T) {
	tests := []struct {
		name   string
		fields map[string][]error
		want   string
	}{
		{
			"single fileld",
			map[string][]error{
				"fieldA": {
					errors.New("too long"),
					errors.New("something wrong"),
				},
			},
			"(fieldA) too long, something wrong",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verr := &ValidationError{fieldErrors: tt.fields}
			got := verr.Error()
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("verror mismatch (-want, +got):%s\n", diff)
			}
		})
	}
}
