package dependency

import (
	"go/build"
	"log"
	"path/filepath"
	"strings"
	"sync"
)

func Walk(ctx build.Context, projectDir string, walkFunc func(pkg *build.Package) error) (err error) {
	errChan := make(chan error, 0)
	done := false

	go func() {
		defer close(errChan)
		wg := sync.WaitGroup{}
		mapProcessImports := sync.Map{}

		if rootPkg, importDirErr := ctx.ImportDir(projectDir, 0); importDirErr != nil {
			log.Println("Failed to import projectDir: ", importDirErr)
			errChan <- importDirErr
			return
		} else {
			wg.Add(1)
			go func() {
				processImports(ctx, rootPkg, walkFunc, &wg, &mapProcessImports, errChan, &done)
				wg.Done()
			}()
			wg.Wait()
		}
	}()

	select {
	case err = <-errChan:
		if err != nil {
			log.Println("Failed to walk: ", err)
			done = true
		}
	}
	return err
}

func processImports(ctx build.Context, pkg *build.Package, walkFunc func(pkg *build.Package) error, wg *sync.WaitGroup, mapProcessImports *sync.Map, errChan chan error, done *bool) {
	value, ok := mapProcessImports.Load(pkg.ImportPath)
	if (ok && value.(bool)) || *done {
		return
	}
	mapProcessImports.Store(pkg.ImportPath, true)
	if err := walkFunc(pkg); err != nil {
		errChan <- err
		return
	}
	for _, pkgName := range pkg.Imports {
		var child *build.Package
		var err error
		if strings.HasPrefix(pkgName, ".") {
			if child, err = ctx.Import(pkgName, pkg.Dir, 0); err != nil {
				log.Println("Failed to Import child start with .: ", err)
				errChan <- err
				return
			}
			child.ImportPath = filepath.Join(pkg.ImportPath, child.ImportPath)
		} else {
			if child, err = ctx.Import(pkgName, "", 0); err != nil {
				log.Println("Failed to Import normal child: ", err)
				errChan <- err
				return
			}
		}
		if child.Goroot {
			continue
		}
		wg.Add(1)
		go func() {
			processImports(ctx, child, walkFunc, wg, mapProcessImports, errChan, done)
			wg.Done()
		}()
	}
}
