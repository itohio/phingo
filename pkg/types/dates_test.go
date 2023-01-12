package types

import "testing"

func TestSanitizeDateTime(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		{"2022-12-31 20:21", "2022-12-31 20:21", false},
		{"2022/12/31 20:21", "2022-12-31 20:21", false},
		{"20221231 20:21", "2022-12-31 20:21", false},
		{"2022-12-31T20:21", "2022-12-31 20:21", false},
		{"2022/12/31T20:21", "2022-12-31 20:21", false},
		{"20221231T20:21", "2022-12-31 20:21", false},
		{"2022-12-31", "2022-12-31 00:00", false},
		{"2022/12/31", "2022-12-31 00:00", false},
		{"20221231", "2022-12-31 00:00", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SanitizeDateTime(tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("SanitizeDateTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SanitizeDateTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
