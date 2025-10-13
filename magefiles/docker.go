package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Docker namespace for docker-related commands
type Docker mg.Namespace

// DevUp starts development Docker services
func (p Docker) Up() error {
	fmt.Println("Starting development Docker services...")
	return sh.RunV("docker-compose", "-f", "docker/docker-compose.dev.yml", "up", "-d")
}

// DevDown stops development Docker services
func (p Docker) Down() error {
	fmt.Println("Stopping development Docker services...")
	return sh.RunV("docker-compose", "-f", "docker/docker-compose.dev.yml", "down")
}

func (p Docker) Clean() error {
	fmt.Println("Cleaning development Docker services...")
	return sh.RunV("docker-compose", "-f", "docker/docker-compose.dev.yml", "down", "-v", "--remove-orphans")
}

// Logs shows Docker logs
func (p Docker) Logs() error {
	fmt.Println("Showing Docker logs...")
	return sh.RunV("docker-compose", "-f", "docker/docker-compose.dev.yml", "logs", "-f")
}

// Build builds the Docker image
func (p Docker) Build() error {
	fmt.Println("Building Docker image...")
	return sh.RunV("docker-compose", "-f", "docker/docker-compose.dev.yml", "build")
}

// Rebuild rebuilds Docker images without cache
func (p Docker) Rebuild() error {
	fmt.Println("Rebuilding Docker images (no cache)...")
	return sh.RunV("docker-compose", "-f", "docker/docker-compose.dev.yml", "build", "--no-cache")
}

func (p Docker) StartClean() error {
	mg.SerialDeps(p.Clean, p.Rebuild, p.Up)
	fmt.Println("Restarting successfull...")
	return nil
}
