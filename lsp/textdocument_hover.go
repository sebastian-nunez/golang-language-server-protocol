package lsp

func NewTextDocumentHoverResponse(id int, contents MarkedString) *TextDocumentHoverResponse {
	return &TextDocumentHoverResponse{
		Response: Response{
			RPC: "2.0",
			ID:  id,
		},
		Result: &HoverResult{
			Contents: contents,
		},
	}
}

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
