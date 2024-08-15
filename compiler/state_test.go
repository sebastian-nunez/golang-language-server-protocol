package compiler

import (
	"errors"
	"reflect"
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
			wantText:    "package main\n\nfunc main() {}\n", // No change want
			wantErr:     errors.New("document was already opened"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			state := NewState()
			if tc.initialURI != "" {
				state.documents[tc.initialURI] = tc.initialText
			}
			_, err := state.OpenDocument(tc.newURI, tc.newText)

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
			_, err := state.UpdateDocument(tc.updateURI, tc.updateText)

			if err != nil && err.Error() != tc.wantErr.Error() {
				t.Errorf("UpdateDocument got error = %v, want %v", err, tc.wantErr)
			}

			if gotText := state.documents[tc.updateURI]; gotText != tc.wantText {
				t.Errorf("UpdateDocument got text = %v, want %v", gotText, tc.wantText)
			}
		})
	}
}

func TestHover(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		documents   map[lsp.DocumentURI]string
		uri         lsp.DocumentURI
		id          int
		position    lsp.Position
		wantContent lsp.MarkedString
		wantErr     error
	}{
		{
			name:        "existing document",
			documents:   map[lsp.DocumentURI]string{"file:///example.go": "package main\n\nfunc main() {}\n"},
			uri:         lsp.DocumentURI("file:///example.go"),
			id:          1,
			position:    lsp.Position{Line: 1, Character: 5},
			wantContent: lsp.MarkedString("file=file:///example.go, characters=29"),
			wantErr:     nil,
		},
		{
			name:        "non-existing document",
			documents:   map[lsp.DocumentURI]string{},
			uri:         lsp.DocumentURI("file:///nonexistent.go"),
			id:          2,
			position:    lsp.Position{Line: 1, Character: 5},
			wantContent: "",
			wantErr:     ErrDocumentNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			state := &State{documents: tc.documents}
			got, err := state.Hover(tc.uri, tc.id, tc.position)

			if err != nil && err.Error() != tc.wantErr.Error() {
				t.Errorf("Hover got error = %v, want = %v", err, tc.wantErr)
			}

			if got != nil && got.Result.Contents != tc.wantContent {
				t.Errorf("Hover got content = %v, want = %v", got.Result.Contents, tc.wantContent)
			}
		})
	}
}

func TestDefinition(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		documents map[lsp.DocumentURI]string
		uri       lsp.DocumentURI
		id        int
		position  lsp.Position
		wantRange lsp.Range
		wantURI   lsp.DocumentURI
		wantErr   error
	}{
		{
			name:      "existing document",
			documents: map[lsp.DocumentURI]string{"file:///example.go": "package main\n\nfunc main() {}\n"},
			uri:       lsp.DocumentURI("file:///example.go"),
			id:        1,
			position:  lsp.Position{Line: 2, Character: 10},
			wantRange: lsp.Range{Start: lsp.Position{Line: 1, Character: 0}, End: lsp.Position{Line: 1, Character: 0}},
			wantURI:   lsp.DocumentURI("file:///example.go"),
			wantErr:   nil,
		},
		{
			name:      "non-existing document",
			documents: map[lsp.DocumentURI]string{},
			uri:       lsp.DocumentURI("file:///nonexistent.go"),
			id:        1,
			position:  lsp.Position{Line: 2, Character: 10},
			wantRange: lsp.Range{},
			wantURI:   lsp.DocumentURI("file:///nonexistent.go"),
			wantErr:   ErrDocumentNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			state := &State{documents: tc.documents}
			got, err := state.Definition(tc.uri, tc.id, tc.position)

			if err != nil && err.Error() != tc.wantErr.Error() {
				t.Errorf("Definition got error = %v, want %v", err, tc.wantErr)
			}

			if got != nil {
				if *got.Result.Range != tc.wantRange {
					t.Errorf("Definition got range = %v, want %v", *got.Result.Range, tc.wantRange)
				}
				if got.Result.URI != tc.wantURI {
					t.Errorf("Definition got uri = %v, want %v", got.Result.URI, tc.wantURI)
				}
			}
		})
	}
}

func TestTextDocumentCodeAction(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		documents map[lsp.DocumentURI]string
		id        int
		uri       lsp.DocumentURI
		want      lsp.TextDocumentCodeActionResponse
		wantError error
	}{
		{
			name: "Document contains 'VS Code'",
			documents: map[lsp.DocumentURI]string{
				"file:///example": "This is a line with VS Code",
			},
			id:  1,
			uri: "file:///example",
			want: lsp.TextDocumentCodeActionResponse{
				Response: lsp.Response{
					RPC: "2.0",
					ID:  1,
				},
				Result: []lsp.CodeAction{
					{
						Title: "Replace VS C*de with a superior editor",
						Edit: &lsp.WorkspaceEdit{
							Changes: map[string][]lsp.TextEdit{
								"file:///example": {
									{
										Range:   LineRange(0, 23, 30), // Adjust based on actual position
										NewText: "Neovim",
									},
								},
							},
						},
					},
					{
						Title: "Censor to VS C*de",
						Edit: &lsp.WorkspaceEdit{
							Changes: map[string][]lsp.TextEdit{
								"file:///example": {
									{
										Range:   LineRange(0, 23, 30), // Adjust based on actual position
										NewText: "VS C*de",
									},
								},
							},
						},
					},
				},
			},
			wantError: nil,
		},
		{
			name: "Document does not contain 'VS Code'",
			documents: map[lsp.DocumentURI]string{
				"file:///example": "No special text here",
			},
			id:  1,
			uri: "file:///example",
			want: lsp.TextDocumentCodeActionResponse{
				Response: lsp.Response{
					RPC: "2.0",
					ID:  1,
				},
				Result: []lsp.CodeAction{}, // No actions should be generated
			},
			wantError: nil,
		},
		{
			name:      "Document not found",
			documents: map[lsp.DocumentURI]string{},
			id:        1,
			uri:       "file:///missing",
			want:      lsp.TextDocumentCodeActionResponse{},
			wantError: ErrDocumentNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			state := &State{documents: tc.documents}
			response, err := state.TextDocumentCodeAction(tc.id, tc.uri)
			if err != nil && err != tc.wantError {
				t.Errorf("want error %v, got %v", tc.wantError, err)
			}

			if err == nil && len(response.Result) == 1 && response.Result[0].Title != tc.want.Result[0].Title {
				t.Errorf("want Title %+v, got %+v", tc.want.Result, response.Result)
			}
		})
	}
}

func TestLineRange(t *testing.T) {
	testCases := []struct {
		name   string
		line   int
		start  int
		end    int
		expect lsp.Range
	}{
		{
			name:  "Standard case",
			line:  1,
			start: 2,
			end:   5,
			expect: lsp.Range{
				Start: lsp.Position{
					Line:      1,
					Character: 2,
				},
				End: lsp.Position{
					Line:      1,
					Character: 5,
				},
			},
		},
		{
			name:  "Start and end are the same",
			line:  2,
			start: 3,
			end:   3,
			expect: lsp.Range{
				Start: lsp.Position{
					Line:      2,
					Character: 3,
				},
				End: lsp.Position{
					Line:      2,
					Character: 3,
				},
			},
		},
		{
			name:  "Zero values",
			line:  0,
			start: 0,
			end:   0,
			expect: lsp.Range{
				Start: lsp.Position{
					Line:      0,
					Character: 0,
				},
				End: lsp.Position{
					Line:      0,
					Character: 0,
				},
			},
		},
		{
			name:  "Negative values",
			line:  -1,
			start: -2,
			end:   -1,
			expect: lsp.Range{
				Start: lsp.Position{
					Line:      -1,
					Character: -2,
				},
				End: lsp.Position{
					Line:      -1,
					Character: -1,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := LineRange(tc.line, tc.start, tc.end)
			if result != tc.expect {
				t.Errorf("want %+v, got %+v", tc.expect, result)
			}
		})
	}
}

func TestTextDocumentCompletion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		id   int
		uri  lsp.DocumentURI
		want *lsp.TextDocumentCompletionResponse
	}{
		{
			name: "Basic Completion",
			id:   1,
			uri:  "file://testfile.go",
			want: &lsp.TextDocumentCompletionResponse{
				Response: lsp.Response{
					RPC: "2.0",
					ID:  1,
				},
				Result: []lsp.CompletionItem{
					{
						Label:         "Custom completion",
						Detail:        "Some super great details.",
						Documentation: "This is a documentation tooltip. In a real app, this would be useful information.",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := &State{}

			got := s.TextDocumentCompletion(tt.id, tt.uri)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TextDocumentCompletion got = %v, want %v", got, tt.want)
			}
		})
	}
}
