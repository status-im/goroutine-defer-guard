package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/status-im/goroutine-defer-guard/pkg/analyzer"
)

func main() {
	a := analyzer.New(nil)

	// singlechecker runs the analyzer
	singlechecker.Main(a)
}
