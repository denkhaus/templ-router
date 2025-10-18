package main

import (
	"fmt"
	"sync"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Default target to run when none is specified
var Default = Dev

// Dev starts the development environment with parallel watch processes
func Dev() error {
	mg.Deps(Build.TailwindClean, Build.TemplGenerate, Build.RegistryGenerate)

	fmt.Println("Starting development server...")

	var wg sync.WaitGroup
	errChan := make(chan error, 4)

	// Start Tailwind watch
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Starting Tailwind CSS watch...")
		if err := (Build{}).TailwindWatch(); err != nil {
			errChan <- fmt.Errorf("tailwind watch failed: %w", err)
		}
	}()

	// Start Templ watch
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Starting Templ watch...")
		if err := (Build{}).TemplWatch(); err != nil {
			errChan <- fmt.Errorf("templ watch failed: %w", err)
		}
	}()

	// Start Registry watch
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Starting Registry watch...")
		if err := (Build{}).RegistryWatch(); err != nil {
			errChan <- fmt.Errorf("registry watch failed: %w", err)
		}
	}()

	// Start Air server
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Starting Air development server...")
		if err := sh.RunV("air",
			"--build.cmd", "cd ./demo && go build -o /tmp/templ-router-demo/main ./main.go",
			"--build.bin", "/tmp/templ-router-demo/main",
			"--build.delay", "100",
			"--build.exclude_dir", "**/node_modules/**",
			"--build.include_ext", "go", "yaml",
			"--build.stop_on_error", "false",
			"--misc.clean_on_exit", "true"); err != nil {
			errChan <- fmt.Errorf("air server failed: %w", err)
		}
	}()

	// Wait for first error or completion
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Return first error if any
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// Clean removes build artifacts and temporary files
func Clean() error {
	fmt.Println("Cleaning build artifacts...")

	// Remove build artifacts
	dirs := []string{
		"/tmp/templui-starter",
		"dist",
		"node_modules/.cache",
	}

	for _, dir := range dirs {
		if err := sh.Rm(dir); err != nil {
			fmt.Printf("Warning: Could not remove %s: %v\n", dir, err)
		}
	}

	return nil
}

// // checkEnvFile checks if .env file exists and creates it from example
// func checkEnvFile() error {
// 	if _, err := os.Stat(".env"); os.IsNotExist(err) {
// 		fmt.Println("Creating .env file from .env.example...")
// 		return sh.Copy(".env", ".env.example")
// 	}
// 	return nil
// }

// // init runs when mage is imported
// func init() {
// 	// Ensure .env file exists
// 	if err := checkEnvFile(); err != nil {
// 		fmt.Printf("Warning: Could not create .env file: %v\n", err)
// 	}
// }
