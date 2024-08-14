package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/sebastian-nunez/golang-language-server-protocol/compiler"
	"github.com/sebastian-nunez/golang-language-server-protocol/lsp"
	"github.com/sebastian-nunez/golang-language-server-protocol/rpc"
	"github.com/sebastian-nunez/golang-language-server-protocol/util"
)

func main() {
	// Since the LSP will be using `os.Stdout` to send messages,
	// we are unable to use `os.Stdout` to log messages.
	logger := util.NewFileLogger("lsp_logs.txt")
	logger.Println("Starting the LSP...")

	state := compiler.NewState()
	writer := os.Stdout

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.SplitMessage)
	for scanner.Scan() {
		msg := scanner.Bytes()
		method, content, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Error decoding message: %v", err)
			continue
		}
		handleMessage(logger, state, writer, method, content)
	}
}

// handleMessage handles the incoming message from the client and sends the appropriate response (if needed).
func handleMessage(logger *log.Logger, state *compiler.State, writer io.Writer, method string, content []byte) {
	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("Error unmarshalling initialize request: %v", err)
			return
		}

		version := "unknown"
		if request.Params.ClientInfo != nil && *request.Params.ClientInfo.Version != "" {
			version = *request.Params.ClientInfo.Version
		}
		logger.Printf("Connected to client: %s (version=%s)",
			request.Params.ClientInfo.Name,
			version,
		)

		response := lsp.NewInitializeResponse(request.ID)
		writeResponse(writer, response)
		logger.Println("Sent initialize response")
	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("Error unmarshalling textDocument/didOpen: %v", err)
			return
		}

		logger.Printf("Opened text document: URI=%v", request.Params.TextDocument.URI)
		err := state.OpenDocument(request.Params.TextDocument.URI, request.Params.TextDocument.Text)
		if err != nil {
			logger.Printf("Error opening document: %v", err)
		}
	case "textDocument/didChange":
		var request lsp.TextDocumentDidChangeNotification
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("Error unmarshalling textDocument/didChange: %v", err)
			return
		}

		for _, change := range request.Params.ContentChanges {
			logger.Printf("Changed text document: URI=%v", request.Params.TextDocument.URI)
			err := state.UpdateDocument(request.Params.TextDocument.URI, change.Text)
			if err != nil {
				logger.Printf("Error updating document: %v", err)
			}
		}
	case "textDocument/hover":
		var request lsp.TextDocumentHoverRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("Error unmarshalling textDocument/hover: %v", err)
			return
		}

		logger.Printf("Hovered over text document: URI=%v, character=%v, line=%v",
			request.Params.TextDocument.URI,
			request.Params.Position.Character,
			request.Params.Position.Line,
		)

		response, err := state.Hover(request.Params.TextDocument.URI, request.ID, request.Params.Position)
		if err != nil {
			logger.Printf("Error getting hover response: %v", err)
		}

		writeResponse(writer, response)
		logger.Println("Sent hover response")
	case "textDocument/definition":
		var request lsp.TextDocumentDefinitionRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("Error unmarshalling textDocument/definition: %v", err)
			return
		}

		logger.Printf("Definition of text document: URI=%v, character=%v, line=%v",
			request.Params.TextDocument.URI,
			request.Params.Position.Character,
			request.Params.Position.Line,
		)

		response, err := state.Definition(request.Params.TextDocument.URI, request.ID, request.Params.Position)
		if err != nil {
			logger.Printf("Error getting definition response: %v", err)
		}

		writeResponse(writer, response)
		logger.Println("Sent definition response")

	case "textDocument/codeAction":
		var request lsp.TextDocumentDefinitionRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("Error unmarshalling textDocument/codeAction: %v", err)
			return
		}

		response, err := state.TextDocumentCodeAction(request.ID, request.Params.TextDocument.URI)
		if err != nil {
			logger.Printf("Error getting codeAction response: %v", err)
		}

		writeResponse(writer, response)
		logger.Println("Sent codeAction response")
	case "textDocument/completion":
		var request lsp.TextDocumentCompletionRequest
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("Error unmarshalling textDocument/completion: %v", err)
			return
		}

		response, err := state.TextDocumentCompletion(request.ID, request.Params.TextDocument.URI)
		if err != nil {
			logger.Printf("Error getting completion response: %v", err)
		}

		writeResponse(writer, response)
		logger.Println("Sent completion response")
	default:
		logger.Printf("Received message: method=%v, content=%v", method, string(content))
	}
}

func writeResponse(writer io.Writer, msg any) {
	encodedMsg := rpc.EncodeMessage(msg)
	writer.Write([]byte(encodedMsg))
}
