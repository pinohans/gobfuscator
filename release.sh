GOOS=windows GOARCH=amd64 go build -o obfuscator.1.0.0.windows.amd64.exe .
GOOS=windows GOARCH=386 go build -o obfuscator.1.0.0.windows.386.exe .
GOOS=linux GOARCH=amd64 go build -o obfuscator.1.0.0.linux.amd64 .
GOOS=linux GOARCH=386 go build -o obfuscator.1.0.0.linux.386 .
GOOS=darwin GOARCH=amd64 go build -o obfuscator.1.0.0.darwin.amd64 .
