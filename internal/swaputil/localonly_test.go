package swaputil

import (
	"net/http"
	"testing"
)

func TestLocalOnlyRequested(t *testing.T) {
	tests := []struct {
		name   string
		absent bool
		vals   []string
		want   bool
	}{
		{name: "absent", absent: true, want: false},
		{name: "empty value counts as present", vals: []string{""}, want: true},
		{name: "1", vals: []string{"1"}, want: true},
		{name: "true", vals: []string{"true"}, want: true},
		{name: "TRUE", vals: []string{"TRUE"}, want: true},
		{name: "padded true", vals: []string{"  true  "}, want: true},
		{name: "0", vals: []string{"0"}, want: false},
		{name: "false", vals: []string{"false"}, want: false},
		{name: "False", vals: []string{"False"}, want: false},
		{name: "padded false", vals: []string{"  false "}, want: false},
		{name: "unknown word counts as present", vals: []string{"maybe"}, want: true},
		{name: "no counts as present", vals: []string{"no"}, want: true},
		{name: "any non-false value among repeats", vals: []string{"false", "1"}, want: true},
		{name: "repeats that are all false", vals: []string{"0", "false"}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var h http.Header
			if !tt.absent {
				h = make(http.Header)
				for _, v := range tt.vals {
					h.Add(HeaderLocalOnly, v)
				}
			}
			if got := LocalOnlyRequested(h); got != tt.want {
				t.Fatalf("LocalOnlyRequested(%v) = %v, want %v", tt.vals, got, tt.want)
			}
		})
	}
}
