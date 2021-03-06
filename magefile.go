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

var devMode bool

// Dev runs all prebuild setps and runs cmd/gopoll in ./run
func Dev() error {
	devMode = true
	prebuild()
	_ = os.Mkdir("run", os.ModeDir)
	os.Chdir("run")
	defer func() {
		os.Chdir("..")
	}()
	fmt.Println("Starting github.com/codemicro/gopoll/cmd/gopoll")
	return sh.RunV("go", "run", "github.com/codemicro/gopoll/cmd/gopoll")
}

// InstallDeps installs dependencies for the Go portion of the project
func InstallDeps() error {
	return sh.Run("go", "mod", "download")
}

func prebuild() {
	mg.Deps(InstallDeps)
	mg.Deps(GenerateTemplates)
	mg.Deps(BundleResources)
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
					fmt.Printf("Moving %s to %s\n", filename, newFilename)
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

// BundleResources takes various web resources and generates the ./internal/webRes package
func BundleResources() error {
	mg.Deps(NPM.BuildStyles)

	// ensure go-bindata is available
	if err := exsh.EnsureGoBin("go-bindata", "github.com/go-bindata/go-bindata/..."); err != nil {
		return err
	}

	command := []string{
		"--nometadata", "-o", alib.OsPathJoin("internal", "webRes", "webRes.go"), "--pkg", "webRes", "--prefix", "build",
	}

	filesToBundle := []string{outputCSSFilename}

	// go-bindata --nometadata -o internal/webRes/webRes.go --pkg webRes --prefix build <files>
	return sh.Run("go-bindata", append(command, filesToBundle...)...)
}

type NPM mg.Namespace

// ./web Node stuff

// InstallDeps installs NPM dependencies for ./web
func (NPM) InstallDeps() error {

	if !exsh.IsCmdAvail("npm") || !exsh.IsCmdAvail("npx") {
		return errors.New("npm and/or npx cannot be found on PATH - see https://nodejs.org/en/")
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

	forceBuildStyles bool
)

// ForceBuild forces NPM styles to be built regardless of file modification times
func (NPM) ForceBuild() {
	forceBuildStyles = true
}

// BuildStyles builds Tailwind CSS styles from ./web/css
func (NPM) BuildStyles() error {

	mg.Deps(NPM.InstallDeps)

	_ = os.Mkdir("build", os.ModeDir)

	sourcesToCheck := []string{
		// These sources should always trigger a recompile
		alib.OsPathJoin("web", "tailwind.config.js"),
		alib.OsPathJoin("web", "css"),
	}

	if !devMode {
		// If we're not in dev mode, changes to the templates should also trigger a CSS rebuild due to the usage of style purges
		sourcesToCheck = append(sourcesToCheck, alib.OsPathJoin("web", "templates"))
	}

	// Check to see if source CSS has been modified since the last modification of: source CSS, Tailwind config file or templates
	if sourcesNewer, err := target.Dir(outputCSSFilename, sourcesToCheck...); err != nil {
		return err
	} else if sourcesNewer || forceBuildStyles {
		os.Chdir("web")
		defer func() {
			os.Chdir("..")
		}()

		env := make(map[string]string)
		if !devMode {
			env["NODE_ENV"] = "production"
		} else if mg.Verbose() {	
			fmt.Println("Dev mode enabled, stylesheet purge will not occur")
		}

		var verboseArg string
		if mg.Verbose() {
			verboseArg = "--verbose"
		}
		
		return sh.RunWith(env, "npx", "postcss", inputCSSFilename, "-o", alib.OsPathJoin("..", outputCSSFilename), verboseArg)
	} else if mg.Verbose() {
		fmt.Println("Skipping building styles, no changes since last build")
	}

	return nil
}
