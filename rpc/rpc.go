package rpc

import (
	"encoding/json"
	"fmt"
)

// EncodeMessage encodes the message into a format that can be sent over the network.
//
// As per the LSP specification, the message should be in the following format:
// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#contentPart
func EncodeMessage(msg any) string {
	content, err := json.Marshal(msg)
	if err != nil {
		// If we are unable to encode the messages, this would be a MAJOR issue.
		// For now, this will NOT be handled gracefully.
		panic(err)
	}

	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
}
