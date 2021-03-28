//+build mage

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/codemicro/alib-go/alib"
	"github.com/codemicro/alib-go/mage/exsh"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

func InstallDeps() error {
	fmt.Println("Installing dependencies")
	if err := sh.Run("go", "mod", "download"); err != nil {
		return err
	}

	return exsh.EnsureGoBin("qtc", "github.com/valyala/quicktemplate/qtc")
}

// QuickTemplate templates

var templateDir = alib.OsPathJoin("web", "templates")
var generatedTemplatesOutputDir = alib.OsPathJoin("internal", "pages", "internal", "templates")

func GenerateTemplates() error {

	// run qtc command
	if err := sh.Run("qtc", "-dir="+templateDir, "-skipLineComments"); err != nil {
		return err
	}

	// move templates from templateDir to the corresponding location in ./internal/pages
	err := filepath.Walk(
		templateDir,
		func(filename string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			if strings.HasSuffix(filename, ".qtpl.go") {
				newFilename := strings.Replace(filename, templateDir, generatedTemplatesOutputDir, 1)

				if mg.Verbose() {
					fmt.Printf("Moving %s to %s", filename, newFilename)
				}

				if err := sh.Copy(newFilename, filename); err != nil {
					return err
				}
				return sh.Rm(filename)
			}

			return nil
		},
	)

	return err
}

type NPM mg.Namespace

// ./web Node stuff

func (NPM) InstallNPMDeps() error {

	if !exsh.IsCmdAvail("npm") || !exsh.IsCmdAvail("npx") {
		return errors.New("npm and/or npm cannot be found on PATH - see https://nodejs.org/en/")
	}

	os.Chdir("web")
	defer func() {
		os.Chdir("..")
	}()

	return sh.Run("npm", "install", "--production=false") // production=false to ensure development dependencies are installed
}

// Tailwind CSS stylesheets

const outputCSSFilename = "main.css"

var inputCSSFilename = alib.OsPathJoin("css", "base.css")

func (NPM) BuildStyles() error {

	mg.Deps(NPM.InstallNPMDeps)

	_ = os.Mkdir("build", os.ModeDir)

	os.Chdir("web")
	defer func() {
		os.Chdir("..")
	}()

	return sh.RunWith(map[string]string{"NODE_ENV": "production"}, "npx", "postcss", inputCSSFilename, "-o", alib.OsPathJoin("..", "build", outputCSSFilename))
}
