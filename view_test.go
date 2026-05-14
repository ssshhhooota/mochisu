package main

import "testing"

func TestRenderMarkdown(t *testing.T) {
	tests := []struct {
		name  string
		src   string
		width int
	}{
		{"empty source", "", 80},
		{"plain text", "hello world", 80},
		{"heading", "# Title", 80},
		{"width zero", "hello", 0},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("renderMarkdown panicked: %v", r)
				}
			}()
			got := renderMarkdown(tc.src, tc.width)
			if tc.src != "" && got == "" {
				t.Errorf("got empty output for non-empty src %q", tc.src)
			}
		})
	}
}
