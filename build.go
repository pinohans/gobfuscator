package main

import (
	"go/build"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Build() error {
	ctx := build.Default
	ldflags := `-s -w`
	if config.WindowsHide == "1" {
		ldflags += " -H=windowsgui"
	}
	if config.NoStatic == "1" {
		ldflags += ` -extldflags '-static'`
	}

	ctx.Dir = filepath.Join(newGopath, "src")
	arguments := []string{"build", "-trimpath", "-ldflags", ldflags, "-tags", config.Tags, "-o", outputPath, "."}
	cmd := exec.Command("go", arguments...)
	cmd.Env = environBuild
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = filepath.Join(newGopath, "src")
	log.Println("go " + strings.Join(cmd.Args, " "))
	if err := cmd.Run(); err != nil {
		log.Fatal("Failed to run BuildSrc: ", err)
	}
	return nil
}
