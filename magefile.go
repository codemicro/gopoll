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
	"github.com/magefile/mage/target"
)

// InstallDeps installs dependencies for the Go portion of the project
func InstallDeps() error {
	return sh.Run("go", "mod", "download")
}

// QuickTemplate templates

var templateDir = alib.OsPathJoin("web", "templates")
var generatedTemplatesOutputDir = alib.OsPathJoin("internal", "pages", "internal", "templates")

// GenerateTemplates runs the QuickTemplate compiler against HTML templates found in ./web/templates and copies them to ./internal/pages/internal/templates
func GenerateTemplates() error {

	// ensure qtc is available
	if err := exsh.EnsureGoBin("qtc", "github.com/valyala/quicktemplate/qtc"); err != nil {
		return err
	}

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

// InstallDeps installs NPM dependencies for ./web
func (NPM) InstallDeps() error {

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

var (
	outputCSSFilename = alib.OsPathJoin("build", "main.css")
	inputCSSFilename = alib.OsPathJoin("css", "base.css")
)

// BuildStyles builds Tailwind CSS styles from ./web/css
func (NPM) BuildStyles() error {

	mg.Deps(NPM.InstallDeps)

	_ = os.Mkdir("build", os.ModeDir)

	// Check to see if source CSS has been modified since the last built set of CSS
	if sourcesNewer, err := target.Dir(outputCSSFilename, alib.OsPathJoin("web", "css")); err != nil {
		return err
	} else if sourcesNewer {
		os.Chdir("web")
		defer func() {
			os.Chdir("..")
		}()

		return sh.RunWith(map[string]string{"NODE_ENV": "production"}, "npx", "postcss", inputCSSFilename, "-o", alib.OsPathJoin("..", outputCSSFilename))
	} else if mg.Verbose() {
		fmt.Println("Skipping building styles, no changes since last build")
	}

	return nil
}
