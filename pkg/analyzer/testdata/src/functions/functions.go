package functions

import "fmt"

type common struct{}

func (c common) LogOnPanic() {}

var Common = common{}

func goodAnonymous() {
	go func() {
		defer Common.LogOnPanic()
		fmt.Println("Hello, World!")
	}()
}

func badAnonymous() {
	go func() { // want "missing defer call to LogOnPanic"
		fmt.Println("Hello, World!")
	}()
}

func goodFunction() {
	defer Common.LogOnPanic()
	fmt.Println("Hello, World!")
}

func badFunction() {
	fmt.Println("Hello, World!")
}

func testGoodFunction() {
	go goodFunction()
}

func testBadFunction() {
	go badFunction() // want "missing defer call to LogOnPanic"
}

type Example struct{}

func (e Example) goodMethod() {
	defer Common.LogOnPanic()
	fmt.Println("Hello, World!")
}

func (e Example) badMethod() {
	fmt.Println("Hello, World!")
}

func testGoodMethod() {
	e := Example{}
	go e.goodMethod()
}

func testBadMethod() {
	e := Example{}
	go e.badMethod() // want "missing defer call to LogOnPanic"
}

type NestedExample struct {
	example Example
}

func testNestedGoodMethod() {
	e := NestedExample{
		example: Example{},
	}
	go e.example.goodMethod()
}

func testNestedBadMethod() {
	e := NestedExample{
		example: Example{},
	}
	go e.example.badMethod() // want "missing defer call to LogOnPanic"
}

// Interface-based method call tests
type Runner interface {
	Run()
}

type GoodImpl struct{}

func (GoodImpl) Run() {
	defer Common.LogOnPanic()
	fmt.Println("good impl")
}

type BadImpl struct{}

func (BadImpl) Run() {
	// missing defer on purpose
	fmt.Println("bad impl")
}

func testInterfaceCalls(i Runner) {
	// Calling via interface should check all implementations in package.
	// Since BadImpl.Run lacks the defer, the linter should report here.
	go i.Run() // want "missing defer call to LogOnPanic"
}
