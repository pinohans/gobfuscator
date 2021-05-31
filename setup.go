package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func Setup() error {
	if err := GoModTidy(); err != nil {
		log.Println("Failed to GoModTidy: ", err)
		return err
	}

	if err := os.RemoveAll(filepath.Join(newGopath, "src")); err != nil {
		log.Println("Failed to RemoveAll in CopySrc: ", err)
		return err
	}
	return nil
}

func GoModTidy() error {
	arguments := []string{"mod", "tidy"}
	cmd := exec.Command("go", arguments...)
	cmd.Env = environ
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = mainPkgDir
	if err := cmd.Run(); err != nil {
		log.Println("Failed to run GoModTidy: ", err)
		return err
	}
	return nil
}
