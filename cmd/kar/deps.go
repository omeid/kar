package main

import (
	"go/build"
	"strings"
)

func init() {
	build.Default.BuildTags = append(build.Default.BuildTags, "kar")
}

func deps(dir string, imp string) ([]string, error) {
	pkgs := map[string]struct{}{}
	err := scanDeps(pkgs, dir, imp)

	if err != nil {
		return nil, err
	}

	list := []string{}
	for pkg := range pkgs {
		list = append(list, pkg)
	}
	return list, nil
}

func scanDeps(packages map[string]struct{}, dir, imp string) error {

	pkg, err := build.Import(imp, dir, build.ImportComment)
	if err != nil {
		return err
	}

	for _, imp := range pkg.Imports {

		if imp == "C" {
			continue
		}

		// See https://github.com/golang/go/issues/17417
		// catch internal vendoring in net/http since go 1.7
		if strings.HasPrefix(imp, "golang_org/x/") {
			continue
		}

		if _, ok := packages[imp]; ok {
			continue
		}

		packages[imp] = struct{}{}

		err := scanDeps(packages, dir, imp)
		if err != nil {
			return err
		}
	}

	return nil
}
