package main

import (
	"os"

	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/status-im/goroutine-defer-guard/pkg/analyzer"
	"github.com/status-im/goroutine-defer-guard/pkg/utils"
)

/*
	Set `-skip=<directory>` to skip errors in certain directories.
	If relative, it is relative to the current working directory.
*/

func main() {
	logger := utils.BuildLogger()

	a, err := analyzer.New(logger)
	if err != nil {
		logger.Errorf("failed to create analyzer: %v", err)
		os.Exit(1)
	}

	// singlechecker runs the analyzer
	singlechecker.Main(a)
}
