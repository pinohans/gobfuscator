# gobfuscator

[![gobfuscator](https://github.com/pinohans/gobfuscator/actions/workflows/main.yml/badge.svg)](https://github.com/pinohans/gobfuscator/actions/workflows/main.yml)

Inspired by [gobfuscate](https://github.com/unixpickle/gobfuscate), but gobfuscator is different in totally.


## 1. how to

Just only use it as go.

```bash
gobfuscator build .
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

