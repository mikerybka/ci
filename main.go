package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/mikerybka/util"
)

func main() {
	// If RUN_NOW env var is present, run the build/deploy right away
	if os.Getenv("RUN_NOW") != "" {
		run()
	}

	// If RUN_AT env var is present, run the build every day at RUN_AT time.
	if os.Getenv("RUN_AT") != "" {
		for {
			runAt := os.Getenv("RUN_AT") // format "HH:mm:ss"
			waitUntil(runAt)
			run()
		}
	}
}

func run() {
	// Read config
	configFile := os.Getenv("CONFIG_FILE")
	config := []string{}
	util.ReadJSONFile(configFile, &config)

	// Build and push Docker images
	for _, img := range config {
		fmt.Println("Building", img)
		err := build(img)
		if err != nil {
			fmt.Println("ERROR", err)
		}
	}
}

func build(img string) error {
	err := pull(img)
	if err != nil {
		return err
	}
	return dockerBuildAndPush(img)
}

func pull(img string) error {
	srcDir := os.Getenv("SRC_DIR")
	path := filepath.Join(srcDir, img)

	if !dirExists(path) {
		// Clone if not exists

		// Mkdir
		dir := filepath.Dir(path)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}

		// git clone
		gitURL := fmt.Sprintf("git@github.com:%s.git", img)
		cmd := exec.Command("git", "clone", gitURL)
		cmd.Env = append(os.Environ(), "GIT_SSH_COMMAND=ssh -o StrictHostKeyChecking=no")
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(out))
			return err
		}
	} else {
		// Otherwise pull
		cmd := exec.Command("git", "pull")
		cmd.Dir = path
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(out))
			return err
		}
	}

	return nil
}

func dirExists(dir string) bool {
	fi, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

func dockerBuildAndPush(img string) error {
	srcDir := os.Getenv("SRC_DIR")
	cmd := exec.Command("docker", "buildx", "build",
		"--platform", "linux/amd64,linux/arm64",
		"-t", img,
		"--push",
		".")
	cmd.Dir = filepath.Join(srcDir, img)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func waitUntil(timeOfDay string) {
	// Parse time string
	layout := "15:04:05"
	tt, err := time.Parse(layout, timeOfDay)
	if err != nil {
		panic(err)
	}

	// Calculate duration to wait
	now := time.Now()
	t := time.Date(now.Year(), now.Month(), now.Day(),
		tt.Hour(), tt.Minute(), tt.Second(), tt.Nanosecond(), time.Local)
	if t.Before(now) {
		t = t.Add(24 * time.Hour)
	}
	duration := time.Until(t)

	// Wait
	time.Sleep(duration)
}
