package kjftt

import (
	"testing"
)

func TestKJFTT_GetBookIDFromURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{"empty url", "", ""},
		{"random string", "abcb", ""},
		{"invalid url", "https://example.net", ""},
		{"invalid url query", "https://example.net?idc=2", ""},
		{"valid input", "-id=chamo:11921596&theme=ttkjf", "chamo:11921596"},
		{"valid url", "https://chamo.kis3g.sk/lib/item?id=chamo:11921596", "chamo:11921596"},
		{"valid url with additional param before", "https://chamo.kis3g.sk/lib/item?theme=ttkjf&id=chamo:11921596", "chamo:11921596"},
		{"valid url with additional param after", "https://chamo.kis3g.sk/lib/item?id=chamo:11921596&theme=ttkjf", "chamo:11921596"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kjftt := &KJFTT{}
			if got := kjftt.GetBookIDFromURL(tt.url); got != tt.want {
				t.Errorf("KJFTT.GetBookIDFromURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
