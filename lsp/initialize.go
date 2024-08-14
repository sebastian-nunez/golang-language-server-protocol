package lsp

func NewInitializeResponse(id int) InitializeResponse {
	version := "0.0.0-alpha.0"
	textDocumentSync := TextDocumentSyncKind(1)
	hoverProvider := true
	definitionProvider := true

	return InitializeResponse{
		Response: Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: InitializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync:   &textDocumentSync,
				HoverProvider:      &hoverProvider,
				DefinitionProvider: &definitionProvider,
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
	TextDocumentSync   *TextDocumentSyncKind `json:"textDocumentSync,omitempty"`
	HoverProvider      *bool                 `json:"hoverProvider,omitempty"`
	DefinitionProvider *bool                 `json:"definitionProvider,omitempty"`
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

type DocumentURI string
type TraceValue string

type WorkspaceEdit struct {
	Changes map[string][]TextEdit `json:"changes"`
}

type TextEdit struct {
	Range   Range  `json:"range"`
	NewText string `json:"newText"`
}
