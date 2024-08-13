package lsp

func NewInitializeResponse(id int, textDocumentSync TextDocumentSyncKind) InitializeResponse {
	version := "0.0.0-alpha.0"
	return InitializeResponse{
		Response: Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: InitializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync: &textDocumentSync,
			},
			ServerInfo: &ServerInfo{
				Name:    "golang-lsp",
				Version: &version,
			},
		},
	}
}

type InitializeRequest struct {
	Request
	Params InitializeParams `json:"params"`
}

type InitializeResponse struct {
	Response
	Result InitializeResult `json:"result"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   *ServerInfo        `json:"serverInfo,omitempty"`
}

type ServerCapabilities struct {
	TextDocumentSync *TextDocumentSyncKind `json:"textDocumentSync,omitempty"`
	// Yea, not implementing all of this...
}

// TextDocumentSyncKind represents the type of text document sync. Types are:
// 0: None
// 1: Full
// 2: Incremental
type TextDocumentSyncKind int

type ServerInfo struct {
	Name    string  `json:"name"`
	Version *string `json:"version,omitempty"`
}

type InitializeParams struct {
	ProcessId             *int               `json:"processId"`
	ClientInfo            *ClientInfo        `json:"clientInfo,omitempty"`
	Locale                *string            `json:"locale,omitempty"`
	RootPath              *string            `json:"rootPath,omitempty"`
	RootUri               *string            `json:"rootUri"`
	InitializationOptions interface{}        `json:"initializationOptions,omitempty"`
	Capabilities          ClientCapabilities `json:"capabilities"`
	Trace                 *string            `json:"trace,omitempty"`
	WorkspaceFolders      []WorkspaceFolder  `json:"workspaceFolders,omitempty"`
}

type ClientInfo struct {
	Name    string  `json:"name"`
	Version *string `json:"version,omitempty"`
}

type ClientCapabilities struct {
	// Yea, not implementing all of this...
}

type WorkspaceFolder struct {
	URI  string `json:"uri"`
	Name string `json:"name"`
}

type DocumentUri string
type TraceValue string
