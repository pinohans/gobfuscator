package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/token"
	"gobfuscator/internal/dependency"
	"gobfuscator/internal/filesystem"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func Obfuscate(mainPath string, buildPath string) (mainPkg string, err error) {
	var mapPkgName sync.Map
	var collision sync.Map

	if err = dependency.Walk(build.Default, mainPath, func(pkg *build.Package) error {
		for {
			filename := GetRandomMd5()
			if _, isLoad := collision.LoadOrStore(filename, true); isLoad {
				continue
			}
			if _, isLoad := mapPkgName.LoadOrStore(pkg.ImportPath, filename); isLoad {
				return errors.New("same import path")
			}
			return nil
		}
	}); err != nil {
		log.Println("Failed to obfuscate package names: ", err)
	} else if err = doObfuscate(&mapPkgName, mainPath, buildPath); err != nil {
		log.Println("Failed to doObfuscate: ", err)
	}

	if mainPkgName, ok := mapPkgName.Load("."); !ok {
		err = errors.New("no main pkg name")
		log.Println("Failed to obfuscate: ", err)
	} else {
		mainPkg = mainPkgName.(string)
	}

	return mainPkg, err
}

func isMainPkg(pkg *build.Package) bool {
	return pkg.Name == "main"
}

type visitor struct {
	mapPkgName *sync.Map
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.ImportSpec:
		oldValue := strings.Trim(n.Path.Value, "\"")
		newValue, ok := v.mapPkgName.Load(oldValue)
		if ok {
			n.Path.Value = fmt.Sprintf("\"%s\"", newValue.(string))
		}
	}
	return v
}

func processComment(file *ast.File, src string, dst string) {
	astCommentGroups := file.Comments
	for _, astCommentGroup := range astCommentGroups {
		for _, astComment := range astCommentGroup.List {
			text := astComment.Text
			text = strings.Trim(text, " ")
			if strings.HasPrefix(text, "//go:embed") {
				text = strings.TrimPrefix(text, "//go:embed")
				text = strings.Trim(text, " ")
				for _, dir := range strings.Split(text, " ") {
					if dir != "" {
						if dir == "*" {
							dir = "."
						} else {
							dir = strings.TrimRight(dir, "*")
						}
						absSrc := filepath.Join(src, dir)
						absDst := filepath.Join(dst, dir)
						isDir, _ := filesystem.IsDir(absSrc)
						if isDir {
							if err := os.MkdirAll(absDst, 0755); err != nil {
								continue
							}
							_ = filesystem.CopyDir(absSrc, absDst)
						}
					}
				}

			} else if strings.HasPrefix(text, "//") {
				text = strings.TrimPrefix(text, "//")
				text = strings.Trim(text, " ")
				if strings.HasPrefix(text, "import ") {
					astComment.Text = "//"
				}
			}
		}
	}

}

func writeGoFile(filename string, node ast.Node, set *token.FileSet) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	if err = format.Node(out, set, node); err != nil {
		return err
	}
	return nil
}

func doObfuscate(mapPkgName *sync.Map, mainPath string, buildPath string) error {
	if err := dependency.Walk(build.Default, mainPath, func(pkg *build.Package) error {
		var newPath string
		if newImportPath, ok := mapPkgName.Load(pkg.ImportPath); !ok {
			log.Println("Failed to doObfuscate in Walk when mapPkgName Load: ", pkg.ImportPath)
			return errors.New("mapPkgName Load error")
		} else {
			newPath = filepath.Join(buildPath, "src", newImportPath.(string))
			log.Println(newImportPath, "\t", pkg.ImportPath)
		}
		if err := os.MkdirAll(newPath, 0755); err != nil {
			log.Println("Failed to MkdirAll newPath in Walk of CopySrc: ", err)
			return err
		}
		set := token.NewFileSet()
		pkgMap, err := parser.ParseDir(set, pkg.Dir, func(info fs.FileInfo) bool {
			for _, name := range pkg.IgnoredGoFiles {
				if info.Name() == name {
					return false
				}
			}
			return true
		}, parser.ParseComments)
		if err != nil {
			log.Println("Failed to parser ParseDir: ", err)
			return err
		}
		for pkgName, astPackage := range pkgMap {
			if !isMainPkg(pkg) && pkgName == "main" {
				continue
			} else if strings.HasSuffix(pkgName, "_test") {
				continue
			}
			ast.Walk(&visitor{mapPkgName: mapPkgName}, astPackage)
			for filename, astFile := range astPackage.Files {
				dst := filepath.Join(newPath, filepath.Base(filename))
				processComment(astFile, pkg.Dir, newPath)
				if err = writeGoFile(dst, astFile, set); err != nil {
					log.Println("Failed to writeGoFile: ", err)
					return err
				}
			}
		}
		srcFiles := [][]string{
			pkg.CgoFiles,
			pkg.CFiles,
			pkg.CXXFiles,
			pkg.MFiles,
			pkg.HFiles,
			pkg.FFiles,
			pkg.SFiles,
			pkg.SwigFiles,
			pkg.SwigCXXFiles,
			pkg.SysoFiles,
		}

		for _, list := range srcFiles {
			for _, file := range list {
				src := filepath.Join(pkg.Dir, file)
				dst := filepath.Join(newPath, file)
				if err = filesystem.CopyFile(src, dst); err != nil {
					log.Println("Failed to copyFile in Walk of CopySrc: ", err)
					return err
				}
			}
		}
		return nil
	}); err != nil {
		log.Println("Failed to obfuscate package names: ", err)
		return err
	}
	return nil
}
