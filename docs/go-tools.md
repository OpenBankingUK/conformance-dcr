# Go Tools

In this document we describe some of the Go tools that are being used in this project. Specifically what's being used in development mode

#### Install Go Tools

We currently use

- [`gofmt`](https://golang.org/cmd/gofmt/)
- [`golangci-lint`](https://github.com/golangci/golangci-lint)
- [`goimports`](golang.org/x/tools/cmd/goimports)

In order to install tools locally you can just run

```sh
make tools
```

This command will download and install the binaries in your `$GOBIN` directory. As a general rule of thumb your `$GOBIN` path should be inside `$PATH` so that you can invoke commands from anywhere in your system.

#### Optional: If you are using VSCode

You can configure linting to run on VSCode directly

**`settings.json`**
```json
...
  "go.lintTool": "golangci-lint",
  "go.lintFlags": [
    "--config=${workspaceFolder}/.golangci.yml",
    "--fast",
  ],
...
```
