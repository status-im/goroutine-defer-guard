package analyzer

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

const DefaultTarget = "HandlePanic"

type Target struct {
	PackagePath string
	FuncName    string
}

func (t Target) String() string {
	return fmt.Sprintf("%v.%v", t.PackagePath, t.FuncName)
}

func (t *Target) Set(s string) error {
	t.PackagePath = ""
	t.FuncName = ""

	if s == "" {
		return errors.New("empty target")
	}

	idx := strings.LastIndex(s, ".")
	if idx == -1 {
		t.FuncName = s
		return nil
	}

	pkg := strings.TrimSpace(s[:idx])
	fn := strings.TrimSpace(s[idx+1:])
	if fn == "" {
		return errors.New("target must include a function name")
	}

	t.PackagePath = pkg
	t.FuncName = fn
	return nil
}
