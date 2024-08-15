package compiler

import (
	"errors"
	"fmt"
	"strings"

	"github.com/sebastian-nunez/golang-language-server-protocol/lsp"
)

var (
	ErrDocumentNotFound      = errors.New("document was not opened")
	ErrDocumentAlreadyOpened = errors.New("document was already opened")
)

type State struct {
	// documents is a map of document URIs (file names) to their text contents.
	documents map[lsp.DocumentURI]string
}

func NewState() *State {
	return &State{
		documents: make(map[lsp.DocumentURI]string),
	}
}

func (s *State) OpenDocument(uri lsp.DocumentURI, text string) ([]lsp.Diagnostic, error) {
	_, ok := s.documents[uri]
	if ok {
		return nil, ErrDocumentAlreadyOpened
	}
	s.documents[uri] = text
	return getDiagnosticsForFile(text), nil
}

func (s *State) UpdateDocument(uri lsp.DocumentURI, text string) ([]lsp.Diagnostic, error) {
	_, ok := s.documents[uri]
	if !ok {
		return nil, ErrDocumentNotFound
	}
	s.documents[uri] = text
	return getDiagnosticsForFile(text), nil
}

func (s *State) Hover(uri lsp.DocumentURI, id int, position lsp.Position) (*lsp.TextDocumentHoverResponse, error) {
	doc, ok := s.documents[uri]
	if !ok {
		return nil, ErrDocumentNotFound
	}

	// This is mocked behavior. In a real implementation, you would want to
	// return the actual hover information for the given position in the document.
	contents := lsp.MarkedString(fmt.Sprintf("file=%s, characters=%d", uri, len(doc)))
	return lsp.NewTextDocumentHoverResponse(id, contents), nil
}

func (s *State) Definition(uri lsp.DocumentURI, id int, position lsp.Position) (*lsp.TextDocumentDefinitionResponse, error) {
	_, ok := s.documents[uri]
	if !ok {
		return nil, ErrDocumentNotFound
	}

	// This is mocked behavior: the "definition" is the line above the word.
	// In a real implementation, you would want to return the actual hover information
	// for the given position in the document.
	contents := lsp.MarkedString("definition")
	r := lsp.Range{
		Start: lsp.Position{Line: position.Line - 1, Character: 0},
		End:   lsp.Position{Line: position.Line - 1, Character: 0},
	}
	return lsp.NewTextDocumentDefinitionResponse(id, uri, r, contents), nil
}

func (s *State) TextDocumentCodeAction(id int, uri lsp.DocumentURI) (lsp.TextDocumentCodeActionResponse, error) {
	text, ok := s.documents[uri]
	if !ok {
		return lsp.TextDocumentCodeActionResponse{}, ErrDocumentNotFound
	}

	actions := []lsp.CodeAction{}
	for row, line := range strings.Split(text, "\n") {
		idx := strings.Index(line, "VS Code")
		if idx >= 0 {
			replaceChange := map[string][]lsp.TextEdit{}
			replaceChange[string(uri)] = []lsp.TextEdit{
				{
					Range:   LineRange(row, idx, idx+len("VS Code")),
					NewText: "Neovim",
				},
			}

			actions = append(actions, lsp.CodeAction{
				Title: "Replace VS C*de with a superior editor",
				Edit:  &lsp.WorkspaceEdit{Changes: replaceChange},
			})

			censorChange := map[string][]lsp.TextEdit{}
			censorChange[string(uri)] = []lsp.TextEdit{
				{
					Range:   LineRange(row, idx, idx+len("VS Code")),
					NewText: "VS C*de",
				},
			}

			actions = append(actions, lsp.CodeAction{
				Title: "Censor to VS C*de",
				Edit:  &lsp.WorkspaceEdit{Changes: censorChange},
			})
		}
	}

	response := lsp.TextDocumentCodeActionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  id,
		},
		Result: actions,
	}

	return response, nil
}

func LineRange(line, start, end int) lsp.Range {
	return lsp.Range{
		Start: lsp.Position{
			Line:      line,
			Character: start,
		},
		End: lsp.Position{
			Line:      line,
			Character: end,
		},
	}
}

func (s *State) TextDocumentCompletion(id int, uri lsp.DocumentURI) *lsp.TextDocumentCompletionResponse {
	// In a real app, we would run static analysis.
	items := []lsp.CompletionItem{
		{
			Label:         "Custom completion",
			Detail:        "Some super great details.",
			Documentation: "This is a documentation tooltip. In a real app, this would be useful information.",
		},
	}

	response := &lsp.TextDocumentCompletionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  id,
		},
		Result: items,
	}
	return response
}

func stringToPtr(s string) *string {
	return &s
}

func getDiagnosticsForFile(text string) []lsp.Diagnostic {
	diagnostics := []lsp.Diagnostic{}
	for row, line := range strings.Split(text, "\n") {
		if strings.Contains(line, "VS Code") {
			idx := strings.Index(line, "VS Code")
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range:    LineRange(row, idx, idx+len("VS Code")),
				Severity: lsp.DiagnosticSeverityError,
				Source:   stringToPtr("Common knowledge"),
				Message:  "Please make sure we use good language!!",
			})
		}

		if strings.Contains(line, "Neovim") {
			idx := strings.Index(line, "Neovim")
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range:    LineRange(row, idx, idx+len("Neovim")),
				Severity: lsp.DiagnosticSeverityHint,
				Source:   stringToPtr("Common Sense"),
				Message:  "Great choice ;)",
			})

		}
	}

	return diagnostics
}
