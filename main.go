package main

import (
	"embed"

	"github.com/nuclide-research/tome/cmd"
	"github.com/nuclide-research/tome/internal/corpus"
)

//go:embed platforms/*.json
var platformFS embed.FS

func main() {
	corpus.Init(platformFS)
	cmd.Execute()
}
