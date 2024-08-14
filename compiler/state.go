package compiler

import (
	"errors"
	"fmt"

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

func (s *State) OpenDocument(uri lsp.DocumentURI, text string) error {
	_, ok := s.documents[uri]
	if ok {
		return ErrDocumentAlreadyOpened
	}
	s.documents[uri] = text
	return nil
}

func (s *State) UpdateDocument(uri lsp.DocumentURI, text string) error {
	_, ok := s.documents[uri]
	if !ok {
		return ErrDocumentNotFound
	}
	s.documents[uri] = text
	return nil
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
