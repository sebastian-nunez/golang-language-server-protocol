package lsp

type TextDocumentHoverRequest struct {
	Request
	Params HoverParams `json:"params"`
}

type HoverParams struct {
	TextDocumentPositionParams
}

type TextDocumentHoverResponse struct {
	Response
	Result *HoverResult `json:"result,omitempty"`
}

type HoverResult struct {
	Contents MarkedString `json:"contents"`
	Range    *Range       `json:"range,omitempty"`
}

type MarkedString string
