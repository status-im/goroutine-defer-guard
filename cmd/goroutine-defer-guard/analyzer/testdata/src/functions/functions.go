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
