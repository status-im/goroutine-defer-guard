package custompattern

import (
	"fmt"

	wrong "custompattern/other"
	"custompattern/right"
)

func goodAnonymous() {
	go func() {
		defer right.HandlePanic()
		fmt.Println("Hello, World!")
	}()
}

func badAnonymous() {
	go func() { // want "missing defer call to custompattern/right.HandlePanic"
		fmt.Println("Hello, World!")
	}()
}

func badWrongPackage() {
	go func() { // want "missing defer call to custompattern/right.HandlePanic"
		defer wrong.HandlePanic()
		fmt.Println("Hello, World!")
	}()
}
