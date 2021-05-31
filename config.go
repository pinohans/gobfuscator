package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

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

var (
	config = Config{
		MainPkgDir:  ".",
		OutputPath:  "dist",
		NewGopath:   "pkg",
		Tags:        "",
		GOOS:        "windows",
		GOARCH:      "amd64",
		CGOENABLED:  "0",
		WindowsHide: "1",
		NoStatic:    "1",
	}
)

var (
	mainPkgDir = func() string {
		var err error
		var ret string
		if ret, err = filepath.Abs(config.MainPkgDir); err != nil {
			log.Fatalln("Failed to get abs of NewGopath: ", err)
		}
		if err = os.MkdirAll(ret, 0755); err != nil {
			log.Fatalln("Failed to MkdirAll NewGopath: ", err)
		}
		return ret
	}()

	newGopath = func() string {
		var err error
		var ret string
		if ret, err = filepath.Abs(config.NewGopath); err != nil {
			log.Fatalln("Failed to get abs of NewGopath: ", err)
		}
		if err = os.MkdirAll(ret, 0755); err != nil {
			log.Fatalln("Failed to MkdirAll NewGopath: ", err)
		}
		return ret
	}()

	outputPath = func() string {
		if ret, err := filepath.Abs(config.OutputPath); err != nil {
			log.Fatalln("Failed to get abs of OutputDir: ", err)
			return ""
		} else {
			return ret
		}
	}()

	randMd5 = func() func() string {
		rand.Seed(time.Now().Unix())
		return func() string {
			buf := make([]byte, 32)
			rand.Read(buf)
			return fmt.Sprintf("%x", md5.Sum(buf))
		}
	}()

	ctxt = func() *build.Context {
		ret := build.Default
		ret.GOPATH = newGopath
		ret.Dir = mainPkgDir
		ret.GOOS = config.GOOS
		ret.GOARCH = config.GOARCH
		return &ret
	}()

	environ = func() []string {
		return append(os.Environ(),
			"GOOS="+ctxt.GOOS,
			"GOARCH="+ctxt.GOARCH,
			"GOROOT="+ctxt.GOROOT,
			"GOPATH="+ctxt.GOPATH,
			"CGO_ENABLED="+config.CGOENABLED,
			"GO111MODULE=auto",
		)
	}()

	environBuild = func() []string {
		return append(os.Environ(),
			"GOOS="+ctxt.GOOS,
			"GOARCH="+ctxt.GOARCH,
			"GOROOT="+ctxt.GOROOT,
			"GOPATH="+ctxt.GOPATH,
			"CGO_ENABLED="+config.CGOENABLED,
			"GO111MODULE=off",
		)
	}()
)

func (c *Config) LoadConfig(configFilename string) error {

	configFile, err := os.OpenFile(configFilename, os.O_RDONLY, 0644)
	if err != nil {
		log.Println("Failed to Open config File: ", err)
		return err
	}
	configJson, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Println("Failed to ReadAll config File: ", err)
		return err
	}
	configTemp := Config{}
	if err = json.Unmarshal(configJson, &configTemp); err != nil {
		log.Println("Failed to Unmarshal config File: ", err)
		return err
	}
	if configTemp.MainPkgDir != "" {
		config.MainPkgDir = configTemp.MainPkgDir
	}
	if configTemp.OutputPath != "" {
		config.OutputPath = configTemp.OutputPath
	}
	if configTemp.NewGopath != "" {
		config.NewGopath = configTemp.NewGopath
	}
	if configTemp.Tags != "" {
		config.Tags = configTemp.Tags
	}
	if configTemp.GOOS != "" {
		config.GOOS = configTemp.GOOS
	}
	if configTemp.GOARCH != "" {
		config.GOARCH = configTemp.GOARCH
	}
	if configTemp.CGOENABLED != "" {
		config.CGOENABLED = configTemp.CGOENABLED
	}
	if configTemp.WindowsHide != "" {
		config.WindowsHide = configTemp.WindowsHide
	}
	if configTemp.NoStatic != "" {
		config.NoStatic = configTemp.NoStatic
	}
	if err = configFile.Close(); err != nil {
		log.Println("Failed to Close config File: ", err)
		return err
	}
	return nil
}
