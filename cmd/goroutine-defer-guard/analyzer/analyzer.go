package analyzer

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"go.uber.org/zap"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"github.com/status-im/goroutine-defer-guard/cmd/goroutine-defer-guard/utils"
)

const Pattern = "LogOnPanic"

type Analyzer struct {
	logger *zap.Logger
	cfg    *Config
}

func New(logger *zap.Logger) (*analysis.Analyzer, error) {
	cfg := Config{}
	flags, err := cfg.ParseFlags()
	if err != nil {
		return nil, err
	}

	logger.Info("creating analyzer")

	processor := newAnalyzer(logger, &cfg)

	analyzer := &analysis.Analyzer{
		Name:     "goroutinedeferguard",
		Doc:      fmt.Sprintf("reports missing defer call to %s", Pattern),
		Flags:    flags,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run: func(pass *analysis.Pass) (interface{}, error) {
			return processor.Run(pass)
		},
	}

	return analyzer, nil
}

func newAnalyzer(logger *zap.Logger, cfg *Config) *Analyzer {
	return &Analyzer{
		logger: logger.Named("processor"),
		cfg:    cfg.WithAbsolutePaths(),
	}
}

func (p *Analyzer) Run(pass *analysis.Pass) (interface{}, error) {
	inspected, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, errors.New("analyzer is not type *inspector.Inspector")
	}

	// Create a nodes filter for goroutines (GoStmt represents a 'go' statement)
	nodeFilter := []ast.Node{
		(*ast.GoStmt)(nil),
	}

	// Inspect go statements
	inspected.Preorder(nodeFilter, func(n ast.Node) {
		p.ProcessNode(pass, n)
	})

	return nil, nil
}

func (p *Analyzer) ProcessNode(pass *analysis.Pass, n ast.Node) {
	goStmt, ok := n.(*ast.GoStmt)
	if !ok {
		panic("unexpected node type")
	}

	switch fun := goStmt.Call.Fun.(type) {
	case *ast.FuncLit: // anonymous function
		pos := pass.Fset.Position(fun.Pos())
		logger := p.logger.With(
			utils.ZapURI(pos.Filename, pos.Line),
			zap.Int("column", pos.Column),
		)

		logger.Debug("found anonymous goroutine")
		if err := p.checkGoroutine(fun.Body); err != nil {
			p.logLinterError(pass, fun.Pos(), fun.Pos(), err)
		}

	case *ast.SelectorExpr: // method call
		pos := pass.Fset.Position(fun.Sel.Pos())
		p.logger.Info("found method call as goroutine",
			zap.String("methodName", fun.Sel.Name),
			utils.ZapURI(pos.Filename, pos.Line),
			zap.Int("column", pos.Column),
		)

		if err := p.checkGoroutineDefinition(pass, fun); err != nil {
			p.logLinterError(pass, goStmt.Pos(), goStmt.Pos(), err)
		}

	case *ast.Ident: // function call
		pos := pass.Fset.Position(fun.Pos())
		p.logger.Info("found function call as goroutine",
			zap.String("functionName", fun.Name),
			utils.ZapURI(pos.Filename, pos.Line),
			zap.Int("column", pos.Column),
		)

		if err := p.checkGoroutineDefinition(pass, fun); err != nil {
			p.logLinterError(pass, goStmt.Pos(), goStmt.Pos(), err)
		}

	default:
		p.logger.Error("unexpected goroutine type",
			zap.String("type", fmt.Sprintf("%T", fun)),
		)
	}
}

func (p *Analyzer) checkGoroutine(body *ast.BlockStmt) error {
	if body == nil {
		p.logger.Warn("missing function body")
		return nil
	}
	if len(body.List) == 0 {
		// empty goroutine is weird, but it is not a linter error
		return nil
	}

	deferStatement, ok := body.List[0].(*ast.DeferStmt)
	if !ok {
		return errors.New("first statement is not defer")
	}

	selectorExpr, ok := deferStatement.Call.Fun.(*ast.SelectorExpr)
	if !ok {
		return errors.New("first statement call is not a selector")
	}

	firstLineFunName := selectorExpr.Sel.Name
	if firstLineFunName != Pattern {
		return errors.Errorf("first statement is not %s", Pattern)
	}

	return nil
}

func (p *Analyzer) checkGoroutineDefinition(pass *analysis.Pass, fun ast.Expr) error {
	var funcName string
	var receiverType string

	// Extract function name and receiver type if it's a method
	switch e := fun.(type) {
	case *ast.Ident:
		funcName = e.Name
	case *ast.SelectorExpr:
		funcName = e.Sel.Name
		// Try to get the receiver type from the selector expression
		if ident, ok := e.X.(*ast.Ident); ok {
			if obj := pass.TypesInfo.ObjectOf(ident); obj != nil {
				if varObj, ok := obj.(*types.Var); ok {
					receiverType = varObj.Type().String()
				}
			}
		}
	default:
		return errors.New("unsupported function expression type")
	}

	// Find the function declaration using AST traversal
	for _, file := range pass.Files {
		var body *ast.BlockStmt

		ast.Inspect(file, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.FuncDecl:
				// Check if this is the function we're looking for
				if node.Name.Name != funcName {
					break
				}

				// If we're looking for a method, check the receiver
				if receiverType != "" && node.Recv != nil && len(node.Recv.List) > 0 {
					// For methods, we need to match the receiver type
					for _, field := range node.Recv.List {
						receiverTypeFromDecl := types.ExprString(field.Type)
						// Handle pointer receivers and named types
						if strings.Contains(receiverType, receiverTypeFromDecl) ||
							strings.Contains(receiverTypeFromDecl, receiverType) {
							body = node.Body
							return false
						}
					}
				} else if receiverType == "" && node.Recv == nil {
					// For regular functions, no receiver should be present
					body = node.Body
					return false
				}
			}
			return true
		})

		if body != nil {
			return p.checkGoroutine(body)
		}
	}

	return errors.New("could not find function body")
}

func (p *Analyzer) logLinterError(pass *analysis.Pass, errPos token.Pos, callPos token.Pos, err error) {
	errPosition := pass.Fset.Position(errPos)
	callPosition := pass.Fset.Position(callPos)

	if p.skip(errPosition.Filename) || p.skip(callPosition.Filename) {
		return
	}

	message := fmt.Sprintf("missing %s()", Pattern)
	p.logger.Warn(message,
		utils.ZapURI(errPosition.Filename, errPosition.Line),
		zap.String("details", err.Error()))

	if callPos == errPos {
		pass.Reportf(errPos, "missing defer call to %s", Pattern)
	} else {
		pass.Reportf(callPos, "missing defer call to %s", Pattern)
	}
}

func (p *Analyzer) skip(filepath string) bool {
	return p.cfg.SkipDir != "" && strings.HasPrefix(filepath, p.cfg.SkipDir)
}
