package main

import (
	"context"

	"github.com/shopware/shopware-cli/cmd"
	"github.com/shopware/shopware-cli/internal/verifier"
)

func main() {
	verifier.TestFoo()
	cmd.Execute(context.Background())
}
