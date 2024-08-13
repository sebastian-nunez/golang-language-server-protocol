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
		msg := scanner.Text()
		handleMessage(logger, msg)
	}
}

func handleMessage(logger *log.Logger, msg any) {
	logger.Printf("Received message: %v", msg)
}
