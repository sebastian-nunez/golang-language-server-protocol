package lsp

// Request is the structure that all LSP requests should follow.
type Request struct {
	RPC    string `json:"jsonrpc"`
	ID     int    `json:"id"`
	Method string `json:"method"`
	// Params will be specified within each of the request types.
}

// Response is the structure that all LSP responses should follow.
type Response struct {
	RPC string `json:"jsonrpc"`
	ID  *int   `json:"id,omitempty"` // Can be nil.
	// Result and errors will be specified within each of the response types.
}

type Notification struct {
	RPC    string `json:"jsonrpc"`
	Method string `json:"method"`
}
