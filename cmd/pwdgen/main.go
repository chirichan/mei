package main

import (
	"context"

	"github.com/spf13/cobra"
)

func main() {
	cobra.CheckErr(NewCLI().ExecuteContext(context.Background()))
}
