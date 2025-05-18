package tpl

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/umefy/godash/logger"
	"github.com/umefy/umex/flags"
	"github.com/umefy/umex/pkg/git"

	"github.com/gookit/color"
)

const (
	goWebAppRepoURL           = "https://github.com/umefy/go-web-app-template.git"
	goWebAppDefaultModuleName = "github.com/umefy/go-web-app-template"
	goWebAppLocalRepoPath     = "/Users/leizhao/Documents/umefy/go-web-app-template"
)

func CreateGoWebApp(flagsModel flags.Model, projectDir string) error {
	projectModuleName := flagsModel.Module
	localMode := flagsModel.LocalMode
	debugMode := flagsModel.DebugMode

	logLevel := slog.LevelInfo
	if debugMode {
		logLevel = slog.LevelDebug
	}
	loggerOpts := logger.NewLoggerOps(false, os.Stdout, logLevel, true, "source", 3)
	logger := logger.New(loggerOpts, nil)

	logger.Debug("Start create go web app", slog.Any("flags", flagsModel))

	if projectModuleName == "" {
		return fmt.Errorf("project module is required")
	}

	currDir, err := os.Getwd()
	if err != nil {
		return err
	}

	tmpDir := filepath.Join(currDir, "tmp")

	logger.Debug("Start cloning template...")
	err = git.CloneRepo(localMode, goWebAppRepoURL, tmpDir, goWebAppLocalRepoPath)
	if err != nil {
		_ = os.RemoveAll(tmpDir)
		return err
	}

	logger.Debug("Start replace go module...")
	err = replaceGoModule(tmpDir, projectModuleName)
	if err != nil {
		_ = os.RemoveAll(tmpDir)
		return err
	}

	logger.Debug("Start update module import...")
	err = updateModuleImport(tmpDir, projectModuleName)
	if err != nil {
		_ = os.RemoveAll(tmpDir)
		return err
	}

	logger.Debug("Start move all files...")
	err = moveAllFiles(projectDir, tmpDir)
	if err != nil {
		_ = os.RemoveAll(tmpDir)
		return err
	}

	fmt.Println("Project created successfully! 🎉")

	projectBaseDir := filepath.Base(projectDir)
	if projectBaseDir != "." {
		fmt.Printf("Please run %s to go to the project directory.\n", color.Green.Render(fmt.Sprintf("cd %s", projectBaseDir)))
		fmt.Printf("Please %s to match with your own env variables and then %s.\n", color.Green.Render("update envrc.example"), color.Green.Render("rename it to .envrc"))
	}
	fmt.Printf("Then please run %s to setup the project.\n", color.Green.Render("./scripts/local_setup.sh"))
	fmt.Printf("Once setup finish, you can run %s to start the project. 🚀\n", color.Green.Render("make"))

	return nil
}

func updateModuleImport(workingDir string, projectModuleName string) error {
	return filepath.WalkDir(workingDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".go") {
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				return err
			}

			updated := false

			ast.Inspect(f, func(n ast.Node) bool {
				if imp, ok := n.(*ast.ImportSpec); ok {
					if strings.Contains(imp.Path.Value, goWebAppDefaultModuleName) {
						imp.Path.Value = strings.ReplaceAll(imp.Path.Value, goWebAppDefaultModuleName, projectModuleName)
						updated = true
					}
				}
				return true
			})

			if updated {
				outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
				if err != nil {
					return err
				}
				//nolint:errcheck
				defer outFile.Close()

				err = printer.Fprint(outFile, fset, f)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func replaceGoModule(workingDir string, projectModuleName string) error {
	// remove existing go.mod and go.sum

	err := os.Remove(filepath.Join(workingDir, "go.mod"))
	if err != nil {
		return err
	}

	err = os.Remove(filepath.Join(workingDir, "go.sum"))
	if err != nil {
		return err
	}

	// create new go.mod
	cmd := exec.Command("go", "mod", "init", projectModuleName)
	cmd.Dir = workingDir
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard

	return cmd.Run()
}

func moveAllFiles(destDir, srcDir string) error {
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return err
		}
	}

	if err := os.CopyFS(destDir, os.DirFS(srcDir)); err != nil {
		return err
	}

	return os.RemoveAll(srcDir)
}
