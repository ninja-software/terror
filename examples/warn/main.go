package main

import (
	"fmt"
	"strings"

	"github.com/ninja-software/terror/v2"
)

// example of warn, with callback

const appVersion = "v1.2.3"

func isModTwo(i int) error {
	if i > 5 {
		return fmt.Errorf("input too big")
	}
	if i%2 != 0 {
		return fmt.Errorf("not mod 2")
	}
	return nil
}

func main() {
	// set globally
	terror.SetVersion(appVersion)

	for i := 0; i < 7; i++ {
		fmt.Println("input", i)
		err := isModTwo(i)
		if err != nil && strings.Contains(err.Error(), "not mod 2") {
			terror.Echo(terror.Warn(err))
		} else if err != nil {
			terror.Echo(terror.Error(err))
		}
		fmt.Println("-------------------")
	}
}
