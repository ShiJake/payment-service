//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Default target to run when none is specified
var Default = Build

// Proto namespace for protobuf-related tasks
type Proto mg.Namespace

// Generate generates protobuf code from .proto files
func (Proto) Generate() error {
	fmt.Println("Generating protobuf code...")

	if err := os.MkdirAll("rpc/fxservice", 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	return sh.RunV("protoc",
		"--proto_path=api/proto",
		"--go_out=rpc/fxservice",
		"--twirp_out=rpc/fxservice",
		"--go_opt=paths=source_relative",
		"--twirp_opt=paths=source_relative",
		"api/proto/fxservice.proto",
	)
}

// Build builds the server binary
func Build() error {
	proto := Proto{}
	mg.Deps(proto.Generate)
	fmt.Println("Building server...")

	if err := os.MkdirAll("bin", 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	return sh.RunV("go", "build", "-o", "bin/server", "./cmd/server")
}

// Clean removes build artifacts
func Clean() error {
	fmt.Println("Cleaning build artifacts...")
	sh.Rm("bin")
	sh.Rm("rpc")
	return nil
}

// Install installs required build tools
func Install() error {
	fmt.Println("Installing build tools...")

	tools := []string{
		"google.golang.org/protobuf/cmd/protoc-gen-go@latest",
		"github.com/twitchtv/twirp/protoc-gen-twirp@latest",
	}

	for _, tool := range tools {
		fmt.Printf("Installing %s...\n", tool)
		if err := sh.RunV("go", "install", tool); err != nil {
			return fmt.Errorf("failed to install %s: %w", tool, err)
		}
	}

	fmt.Println("All tools installed successfully")
	return nil
}

