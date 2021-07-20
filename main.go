package main

import (
	"crypto/md5"
	"fmt"
	"go/build"
	"gobfuscator/internal/filesystem"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	isBuild = len(os.Args) >= 3 && os.Args[1] == "build"

	currentPath = func() string {
		if ret, err := os.Getwd(); err != nil {
			log.Fatalln("Failed to get current path:", err)
		} else {
			return ret
		}
		return ""
	}()

	outputPath = func() string {
		var err error
		ret := "main"
		for index, arg := range os.Args {
			switch arg {
			case "-o":
				if filepath.IsAbs(os.Args[index+1]) {
					return os.Args[index+1]
				} else {
					ret = os.Args[index+1]
				}
			}
		}
		if ret, err = filepath.Abs(filepath.Join(currentPath, ret)); err != nil {
			log.Fatalln("Failed to get output path:", err)
		} else {
			return ret
		}
		return ""
	}()

	mainPath = func() string {
		if isBuild {
			if ret, err := filepath.Abs(os.Args[len(os.Args)-1]); err != nil {
				log.Fatalln("Failed to get abs path:", err)
			} else {
				if isDir, err := filesystem.IsDir(ret); err != nil {
					log.Fatalln("Failed to get dir:", err)
				} else if isDir {
					return ret
				} else {
					return filepath.Dir(ret)
				}
			}
		}
		return ""
	}()

	buildPath = filepath.Join(mainPath, "build")

	buildEnv = func() []string {
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
	}()

	buildCtx = func() build.Context {
		ret := build.Default

		if dir, err := os.Getwd(); err != nil {
			log.Fatalln("Failed to get current dir: ", err)
		} else {
			ret.Dir = dir
		}

		for _, env := range os.Environ() {
			switch true {
			case strings.HasPrefix(env, "GOARCH="):
				ret.GOARCH = strings.TrimPrefix(env, "GOARCH=")
			case strings.HasPrefix(env, "GOOS="):
				ret.GOOS = strings.TrimPrefix(env, "GOOS=")
			case strings.HasPrefix(env, "GOROOT="):
				ret.GOROOT = strings.TrimPrefix(env, "GOROOT=")
			case strings.HasPrefix(env, "GOPATH="):
				ret.GOPATH = strings.TrimPrefix(env, "GOPATH=")
			case strings.HasPrefix(env, "CGO_ENABLED="):
				if CgoEnabled, err := strconv.Atoi(strings.TrimPrefix(env, "CGO_ENABLED=")); err != nil {
					log.Fatalln("Failed to get current dir: ", err)
				} else {
					ret.CgoEnabled = CgoEnabled != 0
				}
			}
		}
		return ret
	}()

	randMd5 = func() func() string {
		rand.Seed(time.Now().Unix())
		return func() string {
			buf := make([]byte, 32)
			rand.Read(buf)
			return fmt.Sprintf("%x", md5.Sum(buf))
		}
	}()
)

func main() {
	if isBuild {
		if err := os.RemoveAll(buildPath); err != nil {
			log.Fatalln("Failed to RemoveAll: ", err)
		}

		if err := os.MkdirAll(buildPath, 0755); err != nil {
			log.Fatalln("Failed to MkdirAll: ", err)
		}

		if err := Obfuscate(); err != nil {
			log.Fatalln("Failed to obfuscate: ", err)
		}

		args := os.Args[1 : len(os.Args)-1]
		args = append(args, ".")
		for index, arg := range args {
			switch arg {
			case "-o":
				args[index+1] = outputPath
			}
		}

		for index, arg := range args {
			switch arg {
			case "-o":
				if !filepath.IsAbs(args[index+1]) {
					args[index+1] = filepath.Join(currentPath, args[index+1])
				}
			}
		}

		cmd := exec.Command("go", args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = filepath.Join(buildPath, "src")
		cmd.Env = buildEnv
		log.Println(strings.Join(cmd.Args, " "))
		if err := cmd.Run(); err != nil {
			log.Fatalln("Failed to run build src:", err)
		}

		if err := os.RemoveAll(buildPath); err != nil {
			log.Fatalln("Failed to clean build src:", err)
		}

	} else {
		cmd := exec.Command("go", os.Args[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		log.Println(strings.Join(cmd.Args, " "))
		if err := cmd.Run(); err != nil {
			log.Fatalln("Failed to run go cmd:", err)
		}
	}
}
