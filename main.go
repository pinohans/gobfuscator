package main

import (
	"crypto/md5"
	"fmt"
	"gobfuscator/internal/filesystem"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func GetRandomMd5() string {
	rand.Seed(time.Now().Unix())
	buf := make([]byte, 32)
	rand.Read(buf)
	return fmt.Sprintf("%x", md5.Sum(buf))
}

func GetBuildEnv(buildPath string) []string {
	ret := make([]string, 0)
	for _, env := range os.Environ() {
		switch true {
		case strings.HasPrefix(env, "GOPATH="):
		case strings.HasPrefix(env, "GO111MODULE="):
		default:
			ret = append(ret, env)
		}
	}
	ret = append(ret, fmt.Sprintf("GOPATH=%s", buildPath))
	ret = append(ret, "GO111MODULE=off")
	return ret
}

func main() {
	var err error
	var mainPath string
	var isDir bool

	if !(len(os.Args) >= 3 && os.Args[1] == "build") {
		log.Fatalln("Please use gobfuscator in build phase.")
	}

	buildPath := fmt.Sprintf(".gobfuscator.%s", GetRandomMd5())

	if err = os.RemoveAll(buildPath); err != nil {
		log.Fatalln("Failed to RemoveAll: ", err)
	}

	if err = os.MkdirAll(buildPath, 0755); err != nil {
		log.Fatalln("Failed to MkdirAll: ", err)
	}

	if mainPath, err = filepath.Abs(os.Args[len(os.Args)-1]); err != nil {
		log.Fatalln("Failed to get abs main path:", err)
	} else if isDir, err = filesystem.IsDir(mainPath); err != nil {
		log.Fatalln("Failed to get dir:", err)
	} else if !isDir {
		mainPath = filepath.Dir(mainPath)
	}

	if err = Obfuscate(mainPath, buildPath); err != nil {
		log.Fatalln("Failed to obfuscate: ", err)
	}

	cmd := exec.Command("go", os.Args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = filepath.Join(buildPath, "src")
	cmd.Env = GetBuildEnv(buildPath)

	log.Println(strings.Join(cmd.Args, " "))

	if err = cmd.Run(); err != nil {
		log.Fatalln("Failed to run build src:", err)
	}

	if err = os.RemoveAll(buildPath); err != nil {
		log.Fatalln("Failed to clean build src:", err)
	}
}
