package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func CloneRepo(localMode bool, templateGitUrl, destDir, localRepoPath string) error {

	if destDir == "." {
		return fmt.Errorf("destination directory cannot be the current directory")
	}

	fmt.Println("Cloning template...")
	if localMode {
		err := cloneRepoDebugMode(destDir, localRepoPath)
		if err != nil {
			return err
		}
	} else {
		cmd := exec.Command("git", "clone", "--depth", "1", "--quiet", templateGitUrl, destDir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	if err := os.RemoveAll(filepath.Join(destDir, ".git")); err != nil {
		return err
	}

	return nil
}

func cloneRepoDebugMode(destDir, localRepoPath string) error {
	if _, err := os.Stat(destDir); os.IsNotExist(err) {
		if err := os.MkdirAll(destDir, 0755); err != nil {
			return err
		}
	}

	return os.CopyFS(destDir, os.DirFS(localRepoPath))
}
