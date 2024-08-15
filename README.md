# Golang LSP

A prototype [Language Server Protocol (LSP)](https://microsoft.github.io/language-server-protocol/overviews/lsp/overview/) with MVP functions written in [Go](https://go.dev/) with the standard library.

## Architecture

The project is organized into several key packages, each responsible for different aspects of the Language Server Protocol (LSP) implementation.

### `/compiler`

Manages the internal state (open, update, and retrieve information about documents) and smart analysis functionality of the documents handled by the LSP.

Supported functions (using Neovim):

- [x] Hover action (press `shift + k`)
- [x] Goto definition (press `g -> d`)
- [x] Code actions (press `SPACE -> c -> a`, must be over "VS Code" text)
- [x] Autocompletion (begin typing `Custom completion` in `insert mode (i)`)
- [x] Diagnostics (open file with the text `VS Code` and `Neovim` somewhere inside)

_This is just a proof of concept, a lot of the functionality is limited and NOT respresentative of a full-fledged LSP._

### `/lsp`

Defines the structures and types required to implement the LSP. This includes requests, responses, and the capabilities of the server.

### `/rpc`

Handles the encoding and decoding of messages sent between the LSP client and server through [Remote Procedure Calls](https://en.wikipedia.org/wiki/Remote_procedure_call) (RPCs).

### `/logs`

This folder is generated after the LSP is running. It will contain all relevant logs regarding messages, actions and responses taken throughout the lifecycle of the program.

## Getting started

> These instructions are specifically for `macOS`. However, with some minor modifications (file path updates), I believe it should also work on `Windows` and/or `Linux`.

1. Install [Neovim](https://neovim.io/) to try out the LSP.

2. Build the project by running `go build /cmd/main.go`
   (assuming project is opened as root)

3. Create a `load_golang_lsp.lua` and add it to your Neovim configuration (`~/.config/nvim/after/plugin/load_golang_lsp.lua`):

   ```lua
   local client = vim.lsp.start_client {
       name = "golang-language-server-protocol",
       cmd = { "/Users/sebastian/golang-language-server-protocol/main" }, -- Update path to Go binary --
   }

   if not client then
       vim.notify("Client was not configured correctly. Make sure the Go Binary has been generated (run `go build /cmd/main.go`). Also, make sure to update `cmd` field of the `vim.lsp.start_client` to have the correct path to the Go binary.")
       return
   end

   vim.api.nvim_create_autocmd("FileType", {
       pattern = "markdown", -- Language supported by the LSP --
       callback = function()
           vim.lsp.buf_attach_client(0, client)
       end,
   })
   ```

4. Open a markdown file using Neovim. (e.g. `nvim README.md`)

_Note: logs are generated into `/logs/`_
