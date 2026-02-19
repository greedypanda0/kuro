package main

import (
	"github.com/greedypanda0/kuro/cli/cmd"
	"github.com/greedypanda0/kuro/cli/internal/logger"
)

func main() {
	logger.Init(true)
	defer logger.Log.Sync()
	cmd.Execute()
}
