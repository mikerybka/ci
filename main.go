package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	buildScriptRepoPath := os.Getenv("BUILD_SCRIPT_REPO_PATH")
	if buildScriptRepoPath == "" {
		buildScriptRepoPath = "/var/lib/src/build"
	}

	for {
		time.Sleep(time.Hour)
		updated, err := pullBuildScript(buildScriptRepoPath)
		if err != nil {
			fmt.Println("pull error:", err)
		}

		if updated {
			fmt.Println("starting build")
			err = build(buildScriptRepoPath)
			if err != nil {
				fmt.Println("build error:", err)
			}
		}
	}
}

func pullBuildScript(buildScriptRepoPath string) (bool, error) {
	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = buildScriptRepoPath
	out, err := cmd.CombinedOutput()
	if err != nil {
		os.Stdout.Write(out)
		return false, err
	}
	if strings.TrimSpace(string(out)) == "Already up to date." {
		return false, nil
	}
	return true, nil
}

func build(buildScriptRepoPath string) error {
	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = buildScriptRepoPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
