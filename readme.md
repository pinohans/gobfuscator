# gobfuscator

[![gobfuscator](https://github.com/pinohans/gobfuscator/actions/workflows/main.yml/badge.svg)](https://github.com/pinohans/gobfuscator/actions/workflows/main.yml)

Inspired by [gobfuscate](https://github.com/unixpickle/gobfuscate), but gobfuscator is different in totally.


## 1. how to

### 1.1. create config.json

> Content of config.json

```json
{
  "MainPkgDir": ".",
  "OutputPath": "dist",
  "NewGopath": "pkg",
  "Tags": "",
  "GOOS": "windows",
  "GOARCH": "amd64",
  "CGOENABLED": "0",
  "WindowsHide": "1",
  "NoStatic": "1"
}
```

> Explanation of parameter

```go
type Config struct {
	MainPkgDir  string // dir of main package
	OutputPath  string // path of dist
	NewGopath   string // dir of new GOPATH
	Tags        string // tags are passed to the go compiler
	GOOS        string // target os
	GOARCH      string // target arch
	CGOENABLED  string // cgo enable
	WindowsHide string // hide windows GUI
	NoStatic    string // no static link
}
```

### 1.2. compile and run

```bash
go mod tidy && go build -o gobfuscator . && ./gobfuscator -c config.json
```

## 2. dependency

> test with Go 1.16.4

## 3. target

- [![gobfuscator-frp](https://github.com/pinohans/gobfuscator/actions/workflows/frp.yml/badge.svg)](https://github.com/pinohans/gobfuscator/actions/workflows/frp.yml)
- [![gobfuscator-fscan](https://github.com/pinohans/gobfuscator/actions/workflows/fscan.yml/badge.svg)](https://github.com/pinohans/gobfuscator/actions/workflows/fscan.yml)

## 4. technical

1. Obfuscate third party package with ast.
2. Process build tags and go:embed.
3. Fast import graph walker.

And so on.

