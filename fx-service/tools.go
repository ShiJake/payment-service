//go:build tools
// +build tools

// Package tools tracks build tool dependencies
package tools

import (
	_ "github.com/magefile/mage"
	_ "github.com/twitchtv/twirp/protoc-gen-twirp"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)

