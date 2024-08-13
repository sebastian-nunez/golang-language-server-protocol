package main

import (
	"bufio"
	"log"
	"os"

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

		handleMessage(logger, method, content)
	}
}

func handleMessage(logger *log.Logger, method string, content []byte) {
	logger.Printf("Received message: method=%v, content=%v", method, string(content))
}
