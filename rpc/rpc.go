package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// BaseMessage has the structure that all RPC messages should follow.
type BaseMessage struct {
	Method string "json:\"method\""
}

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

// DecodeMessage decodes the message from the network into a format that can be used by the application.
//
// As per the LSP specification, the message should be in the following format:
// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#contentPart
func DecodeMessage(msg []byte) (string, []byte, error) {
	sep := []byte{'\r', '\n', '\r', '\n'}
	header, content, found := bytes.Cut(msg, sep)
	if !found {
		return "", nil, errors.New("unable to find separator in message")
	}

	// Header -> `Content-Length: <number>`
	contentLengthBytes := bytes.TrimPrefix(header, []byte("Content-Length: "))
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return "", nil, fmt.Errorf("unable to parse content length from header '%v': %w", header, err)
	}

	var baseMessage BaseMessage
	actualContent := content[:contentLength]
	if err := json.Unmarshal(actualContent, &baseMessage); err != nil {
		return "", nil, fmt.Errorf("unable to parse content from message '%v': %w", content, err)
	}

	return baseMessage.Method, actualContent, nil
}
