# SpaGo

Frontend tool-kit for Gopher

## Feature

- No dependencies on other packages.
- Minimal function set.
- Distributable as static files(etc.: S3 or GitHub-Pages).
- Navigation by Location Hash.
- Mount/Unmount hook.
- Update WASM when reload on browser.
- TinyGo supported.
- Commandline tool included.
- HTML-like DSL translate to Go-Code.

## Install

Library

```shell
go env -w GOOS=js GOARCH=wasm
go get github.com/nobonobo/spago
go env -u GOOS GOARCH
```

Command-line tool

```shell
go get github.com/nobonobo/spago/cmd/spago
```

## Getting Started

[Getting Started](https://github.com/nobonobo/spago/wiki/Getting-Started)
