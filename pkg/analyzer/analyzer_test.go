package analyzer

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/status-im/goroutine-defer-guard/pkg/utils"
)

func TestMethods(t *testing.T) {
	t.Parallel()

	logger := utils.BuildLogger(zap.DebugLevel)

	a, err := New(logger)
	require.NoError(t, err)

	analysistest.Run(t, analysistest.TestData(), a, "functions")
}
