package main

import (
	"cli/cmd"
	"cli/internal/logger"
)

func main() {
	logger.Init(true)
	defer logger.Log.Sync()
	cmd.Execute()
}
