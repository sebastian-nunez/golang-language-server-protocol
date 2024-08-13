package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/sebastian-nunez/golang-language-server-protocol/lsp"
	"github.com/sebastian-nunez/golang-language-server-protocol/rpc"
	"github.com/sebastian-nunez/golang-language-server-protocol/util"
)

func main() {
	// Since the LSP will be using `os.Stdout` to send messages,
	// we are unable to use `os.Stdout` to log messages.
	logger := util.NewFileLogger("lsp_logs.txt")
	logger.Println("Starting the LSP...")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.SplitMessage)
	for scanner.Scan() {
		msg := scanner.Bytes()
		method, content, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("Error decoding message: %v", err)
			continue
		}

		switch method {
		case "initialize":
			var request lsp.InitializeRequest
			if err := json.Unmarshal(content, &request); err != nil {
				logger.Printf("Error unmarshalling initialize request: %v", err)
				continue
			}

			version := "unknown"
			if request.Params.ClientInfo != nil && *request.Params.ClientInfo.Version != "" {
				version = *request.Params.ClientInfo.Version
			}
			logger.Printf("Connected to client: %s (version=%s)",
				request.Params.ClientInfo.Name,
				version,
			)

			// TODO(sebastian-nunez): refactor to just use the writer directly
			writer := os.Stdout
			msg := rpc.EncodeMessage(lsp.NewInitializeResponse(request.ID, 1))
			writer.Write([]byte(msg))

			logger.Printf("Sent initialize response: %v", string(msg))
		case "textDocument/didOpen":
			var request lsp.DidOpenTextDocumentNotification
			if err := json.Unmarshal(content, &request); err != nil {
				logger.Printf("Error unmarshalling text document did open request: %v", err)
				continue
			}

			logger.Printf("Opened text document: URI=%v, content=%v",
				request.Params.TextDocument.URI,
				request.Params.TextDocument.Text,
			)
		default:
			handleMessage(logger, method, content)
		}
	}
}

func handleMessage(logger *log.Logger, method string, content []byte) {
	logger.Printf("Received message: method=%v, content=%v", method, string(content))
}
