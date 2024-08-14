package lsp

func NewTextDocumentDefinitionResponse(id int, uri DocumentURI, rang Range, contents MarkedString) *TextDocumentDefinitionResponse {
	return &TextDocumentDefinitionResponse{
		Response: Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: &Location{
			URI:   uri,
			Range: &rang,
		},
	}
}

type TextDocumentDefinitionRequest struct {
	Request
	Params DefinitionParams `json:"params"`
}

type DefinitionParams struct {
	TextDocumentPositionParams
}

type TextDocumentDefinitionResponse struct {
	Response
	Result *Location `json:"result,omitempty"`
}

type Location struct {
	URI   DocumentURI `json:"uri"`
	Range *Range      `json:"range,omitempty"`
}
