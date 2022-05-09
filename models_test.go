package notes

import "testing"

func TestParseContentID(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    ContentID
		wantErr bool
	}{
		{
			name:    "valid id with prefix",
			arg:     "content:84b4c7eb-69e0-446f-a13b-b02b5a5fa696",
			want:    ContentID("84b4c7eb-69e0-446f-a13b-b02b5a5fa696"),
			wantErr: false,
		},
		{
			name:    "valid id without prefix",
			arg:     "84b4c7eb-69e0-446f-a13b-b02b5a5fa696",
			want:    ContentID("84b4c7eb-69e0-446f-a13b-b02b5a5fa696"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseContentID(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseContentID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseContentID() = %v, want %v", got, tt.want)
			}
		})
	}
}
