package compiler

import (
	"errors"

	"github.com/sebastian-nunez/golang-language-server-protocol/lsp"
)

type State struct {
	// Documents is a map of document URIs (file names) to their text contents.
	Documents map[lsp.DocumentURI]string
}

func NewState() *State {
	return &State{
		Documents: make(map[lsp.DocumentURI]string),
	}
}

func (s *State) OpenDocument(uri lsp.DocumentURI, text string) error {
	_, ok := s.Documents[uri]
	if ok {
		return errors.New("document was already opened")
	}
	s.Documents[uri] = text
	return nil
}

func (s *State) UpdateDocument(uri lsp.DocumentURI, text string) error {
	_, ok := s.Documents[uri]
	if !ok {
		return errors.New("document was not opened")
	}
	s.Documents[uri] = text
	return nil
}
