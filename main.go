package main

import (
	"os"
	"os/exec"
	"time"
)

func main() {
	buildScriptRepoPath := os.Getenv("BUILD_SCRIPT_REPO_PATH")
	if buildScriptRepoPath == "" {
		buildScriptRepoPath = "/var/src/build"
	}

	for {
		err := pullBuildScript(buildScriptRepoPath)
		if err != nil {
			panic(err)
		}

		err = build(buildScriptRepoPath)
		if err != nil {
			panic(err)
		}

		time.Sleep(time.Hour)
	}
}

func pullBuildScript(buildScriptRepoPath string) error {
	cmd := exec.Command("git", "pull")
	cmd.Dir = buildScriptRepoPath
	out, err := cmd.CombinedOutput()
	if err != nil {
		os.Stdout.Write(out)
		return err
	}
	return nil
}

func build(buildScriptRepoPath string) error {
	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = buildScriptRepoPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
