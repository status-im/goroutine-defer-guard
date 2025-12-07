package analyzer

import (
	"log"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestMethods(t *testing.T) {
	t.Parallel()

	logger := log.Default()
	a := New(logger)

	analysistest.Run(t, analysistest.TestData(), a, "functions")
}
