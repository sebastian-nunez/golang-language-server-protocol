# Golang LSP

A [Language Server Protocol (LSP)](https://microsoft.github.io/language-server-protocol/overviews/lsp/overview/) written in [Go](https://go.dev/).

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

## Good to know

- Messages are decoded according to the [LSP specification](https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification) (payloads in [JSON-RPC](https://www.jsonrpc.org/specification)) and logged into a file within `/logs/`.
