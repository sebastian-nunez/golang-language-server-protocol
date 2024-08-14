package compiler

import (
	"errors"
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
			if state.documents == nil {
				t.Errorf("NewState() Documents map is nil")
			}
			if len(state.documents) != 0 {
				t.Errorf("NewState() Documents map should be empty, got %d", len(state.documents))
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
		wantErr     error
	}{
		{
			name:        "new document",
			initialURI:  "",
			initialText: "",
			newURI:      lsp.DocumentURI("file:///example.go"),
			newText:     "package main\n\nfunc main() {}\n",
			wantLength:  1,
			wantText:    "package main\n\nfunc main() {}\n",
			wantErr:     nil,
		},
		{
			name:        "document already exists",
			initialURI:  lsp.DocumentURI("file:///example.go"),
			initialText: "package main\n\nfunc main() {}\n",
			newURI:      lsp.DocumentURI("file:///example.go"),
			newText:     "package main\n\nfunc main() { println(\"Hello, World!\") }\n",
			wantLength:  1,
			wantText:    "package main\n\nfunc main() {}\n", // No change expected
			wantErr:     errors.New("document was already opened"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			state := NewState()
			if tc.initialURI != "" {
				state.documents[tc.initialURI] = tc.initialText
			}
			err := state.OpenDocument(tc.newURI, tc.newText)

			if err != nil && err.Error() != tc.wantErr.Error() {
				t.Errorf("OpenDocument got error = %v, want %v", err, tc.wantErr)
			}

			if gotLength := len(state.documents); gotLength != tc.wantLength {
				t.Errorf("OpenDocument got length = %v, want %v", gotLength, tc.wantLength)
			}

			if gotText := state.documents[tc.newURI]; gotText != tc.wantText {
				t.Errorf("OpenDocument got text = %v, want %v", gotText, tc.wantText)
			}
		})
	}
}

func TestUpdateDocument(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		initialURI  lsp.DocumentURI
		initialText string
		updateURI   lsp.DocumentURI
		updateText  string
		wantText    string
		wantErr     error
	}{
		{
			name:        "update existing document",
			initialURI:  lsp.DocumentURI("file:///example.go"),
			initialText: "package main\n\nfunc main() {}\n",
			updateURI:   lsp.DocumentURI("file:///example.go"),
			updateText:  "package main\n\nfunc main() { println(\"Hello, World!\") }\n",
			wantText:    "package main\n\nfunc main() { println(\"Hello, World!\") }\n",
			wantErr:     nil,
		},
		{
			name:       "update non-existing document",
			initialURI: "",
			updateURI:  lsp.DocumentURI("file:///nonexistent.go"),
			updateText: "package main\n\nfunc main() {}\n",
			wantText:   "",
			wantErr:    errors.New("document was not opened"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			state := NewState()
			if tc.initialURI != "" {
				state.documents[tc.initialURI] = tc.initialText
			}
			err := state.UpdateDocument(tc.updateURI, tc.updateText)

			if err != nil && err.Error() != tc.wantErr.Error() {
				t.Errorf("UpdateDocument got error = %v, want %v", err, tc.wantErr)
			}

			if gotText := state.documents[tc.updateURI]; gotText != tc.wantText {
				t.Errorf("UpdateDocument got text = %v, want %v", gotText, tc.wantText)
			}
		})
	}
}
