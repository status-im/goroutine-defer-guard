package main

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

func goodFunction() {
	defer Common.LogOnPanic()
	fmt.Println("Hello, World!")
}

func testGoodFunction() {
	go goodFunction()
}

type Example struct{}

func (e Example) goodMethod() {
	defer Common.LogOnPanic()
	fmt.Println("Hello, World!")
}

func testGoodMethod() {
	e := Example{}
	go e.goodMethod()
}

func main() {
	goodAnonymous()
	testGoodFunction()
	testGoodMethod()
}
