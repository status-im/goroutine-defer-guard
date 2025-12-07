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

func TestCustomTarget(t *testing.T) {
	t.Parallel()

	logger := log.Default()
	a := New(logger)
	if err := a.Flags.Set("target", "custompattern/right.HandlePanic"); err != nil {
		t.Fatalf("set Target flag: %v", err)
	}

	analysistest.Run(t, analysistest.TestData(), a, "custompattern")
}
