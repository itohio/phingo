package types

import "testing"

func TestSanitizePath(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"abcdABCD.jkl", "abcdabcd.jkl"},
		{"a~b?c%d$A*B'C\"D.jk`l", "a~b-c-d-a-b-c-d.jk-l"},
		{"abcd\\ABCD.jkl", "abcd-abcd.jkl"},
		{"abcd/ABCD.jkl", "abcd-abcd.jkl"},
		{"ab/../cdABCD.jkl", "ab-cdabcd.jkl"},
		{"../abcdABCD.jkl", "-abcdabcd.jkl"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SanitizePath(tt.name); got != tt.want {
				t.Errorf("SanitizePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
