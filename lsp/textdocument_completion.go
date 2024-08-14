package lsp

func NewTextDocumentCompletionResponse(id int, contents MarkedString) *TextDocumentCompletionResponse {
	return &TextDocumentCompletionResponse{
		Response: Response{
			RPC: "2.0",
			ID:  id,
		},
		Result: []CompletionItem{},
	}
}

type TextDocumentCompletionRequest struct {
	Request
	Params CompletionParams `json:"params"`
}

type CompletionParams struct {
	TextDocumentPositionParams
}

type TextDocumentCompletionResponse struct {
	Response
	Result []CompletionItem `json:"result,omitempty"`
}

type CompletionItem struct {
	Label         string `json:"label"`
	Detail        string `json:"detail,omitempty"`
	Documentation string `json:"documentation,omitempty"`
}
