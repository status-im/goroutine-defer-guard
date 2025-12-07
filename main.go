package main

import (
	"log"
	"os"

	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/status-im/goroutine-defer-guard/pkg/analyzer"
)

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
	a := analyzer.New(logger)

	// singlechecker runs the analyzer
	singlechecker.Main(a)
}
