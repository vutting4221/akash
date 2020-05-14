// +build tools

package tools

//nolint
import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/vektra/mockery"
	_ "k8s.io/code-generator"
)
