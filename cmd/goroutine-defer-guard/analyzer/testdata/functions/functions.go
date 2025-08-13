package functions

import "fmt"

type common struct{}

func (c common) LogOnPanic() {}

var Common = common{}

func goodAnonymous() {
	// This should NOT trigger the linter
	go func() {
		defer Common.LogOnPanic()
		fmt.Println("Hello, World!")
	}()
}

func badAnonymous() {
	// This SHOULD trigger the linter
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
	// This should NOT trigger the linter
	go goodFunction()
}

func testBadFunction() {
	// This SHOULD trigger the linter
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
	// This should NOT trigger the linter
	go e.goodMethod()
}

func testBadMethod() {
	e := Example{}
	// This SHOULD trigger the linter
	go e.badMethod() // want "missing defer call to LogOnPanic"
}
