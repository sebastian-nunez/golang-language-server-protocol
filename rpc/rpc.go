package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

const (
	// ContentLength is the header that specifies the length of the content.
	ContentLength = "Content-Length: "
	// Separator is the separator between the header and the content for an RPC message.
	Separator = "\r\n\r\n"
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

	return fmt.Sprintf("%s%d%v%s", ContentLength, len(content), Separator, content)
}

// DecodeMessage decodes the message from the network into a format that can be used by the application.
//
// As per the LSP specification, the message should be in the following format:
// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#contentPart
func DecodeMessage(msg []byte) (method string, content []byte, err error) {
	sep := []byte(Separator)
	header, content, found := bytes.Cut(msg, sep)
	if !found {
		return "", nil, errors.New("unable to find separator in message: " + string(msg))
	}

	// Header -> `Content-Length: <number>`
	contentLengthBytes := bytes.TrimPrefix(header, []byte(ContentLength))
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

// SplitMessage splits the message to be read by a `bufio.Scanner`.
func SplitMessage(data []byte, _ bool) (advance int, token []byte, err error) {
	sep := []byte(Separator)
	header, content, found := bytes.Cut(data, sep)
	if !found {
		// We are still waiting for more data to come in.
		return 0, nil, nil
	}

	// Header -> `Content-Length: <number>`
	contentLengthBytes := bytes.TrimPrefix(header, []byte(ContentLength))
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return 0, nil, fmt.Errorf("unable to parse content length from header '%v': %w", header, err)
	}

	// Data has not been fully received yet.
	if len(content) < contentLength {
		return 0, nil, nil
	}

	totalLength := len(header) + len(sep) + contentLength
	return totalLength, data[:totalLength], nil
}
