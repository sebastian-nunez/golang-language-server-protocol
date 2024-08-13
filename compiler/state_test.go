package compiler

import (
	"testing"

	"github.com/sebastian-nunez/golang-language-server-protocol/lsp"
)

func TestNewState(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
	}{
		{
			name: "new state",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			state := NewState()
			if state == nil {
				t.Errorf("NewState() returned nil")
			}
			if state.Documents == nil {
				t.Errorf("NewState() Documents map is nil")
			}
			if len(state.Documents) != 0 {
				t.Errorf("NewState() Documents map should be empty, got %d", len(state.Documents))
			}
		})
	}
}

func TestOpenDocument(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		initialURI  lsp.DocumentURI
		initialText string
		newURI      lsp.DocumentURI
		newText     string
		wantLength  int
		wantText    string
	}{
		{
			name:        "new document",
			initialURI:  lsp.DocumentURI(""),
			initialText: "",
			newURI:      lsp.DocumentURI("file:///example.go"),
			newText:     "package main\n\nfunc main() {}\n",
			wantLength:  1,
			wantText:    "package main\n\nfunc main() {}\n",
		},
		{
			name:        "overwrite document",
			initialURI:  lsp.DocumentURI("file:///example.go"),
			initialText: "package main\n\nfunc main() {}\n",
			newURI:      lsp.DocumentURI("file:///example.go"),
			newText:     "package main\n\nfunc main() { println(\"Hello, World!\") }\n",
			wantLength:  1,
			wantText:    "package main\n\nfunc main() { println(\"Hello, World!\") }\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			state := NewState()
			if tc.initialURI != "" {
				state.OpenDocument(tc.initialURI, tc.initialText)
			}
			state.OpenDocument(tc.newURI, tc.newText)

			if gotLength := len(state.Documents); gotLength != tc.wantLength {
				t.Errorf("OpenDocument got length = %v, want %v", gotLength, tc.wantLength)
			}

			if gotText := state.Documents[tc.newURI]; gotText != tc.wantText {
				t.Errorf("OpenDocument got text = %v, want %v", gotText, tc.wantText)
			}
		})
	}
}
