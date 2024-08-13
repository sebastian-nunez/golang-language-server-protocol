package compiler

import "github.com/sebastian-nunez/golang-language-server-protocol/lsp"

type State struct {
	// Documents is a map of document URIs (file names) to their text contents.
	Documents map[lsp.DocumentURI]string
}

func NewState() *State {
	return &State{
		Documents: make(map[lsp.DocumentURI]string),
	}
}

func (s *State) OpenDocument(uri lsp.DocumentURI, text string) {
	s.Documents[uri] = text
}
