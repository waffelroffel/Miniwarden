# Miniwarden

A simple readonly Bitwarden client for Windows using Bitwarden CLI.

## Features

- open autofill with `ctrl+shift+a`
- fuzzy search login entries\*

\*Logins with totp or additional master re-prompt are currently not supported.

### Usage

1. `ctrl+shift+a` to open autofill
2. input search term
3. `enter` to autofill first entry or `tab` to switch to list
4. navigate to desired entry and `enter`

---

## Development

### Prerequisites

Install [Go 1.17](https://golang.org/dl/).

Install [Bitwarden CLI](https://bitwarden.com/help/article/cli/#download-and-install) and add it to the PATH.

Get the project ready with

```
git clone https://github.com/waffelroffel/Miniwarden
go get
go get github.com/akavel/rsrc
rsrc -ico icon.ico -manifest miniwarden.manifest -o rsrc.syso
```

### Run

```
go run .
```

### Build

```
go mod tidy
go build -ldflags -H=windowsgui -o bin/Miniwarden.exe .
```
