package main

import (
	"fmt"

	"github.com/ninja-software/terror/v2"
)

// example of panic

const appVersion = "v1.2.3"

func boom() {
	panic("EXPLOOOOOSION!")
}

func main() {
	// set globally
	terror.SetVersion(appVersion)

	defer func() {
		if r := recover(); r != nil {
			err := terror.Panic(fmt.Errorf(r.(string)))
			if err != nil {
				terror.Echo(err)
			}
		}
	}()

	boom()

	fmt.Println("done")
}
