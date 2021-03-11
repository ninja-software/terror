package main

import (
	"fmt"

	"github.com/ninja-software/terror/v2"
)

// example of warn, with callback

const appVersion = "v1.2.3"

func counting(i *int) {
	// resume if panic
	defer func() {
		if rec := recover(); rec != nil {
			terror.Echo(terror.Panic(fmt.Errorf("out of memory"), "all bets are off").KVs("evasive", "maneuver"))
			*i++
			counting(i)
		}
	}()

	for ; *i < 10; *i++ {
		fmt.Println("input", *i)
		err := doCalc(*i)
		if err != nil {
			terror.Echo(terror.Error(err))
		}
		fmt.Println("-------------------")
	}
}

func doCalc(i int) error {
	err := isModTwo(i)
	if err != nil {
		return terror.Error(err).KVs("val", fmt.Sprintf("%d", i)).KVs("math", "fun")
	}
	return nil
}

func isModTwo(i int) error {
	if i > 7 {
		panic("EXPLOOOOSION!")
	}
	if i > 5 {
		return terror.Error(fmt.Errorf("input too big"), "check input").KVs("avoid", "problem")
	}
	if i%2 != 0 {
		return terror.Warn(fmt.Errorf("not mod 2"), "i can haz").KVs("try", "moar")
	}
	return nil
}

func doWarn(meta terror.Meta, err error) {
	fmt.Println("1111111111 doWarn", meta)
	// panic
	// a := []int{}
	// fmt.Println("....", a[10])
}

func doError(meta terror.Meta, err error) {
	fmt.Println("2222222222 doError", meta)
	// panic
	// a := []int{}
	// fmt.Println("....", a[10])
}

func doPanic(meta terror.Meta, err error) {
	fmt.Println("9999999999 doPanic", meta)
	// panic
	a := []int{}
	fmt.Println("....", a[10])
}

func main() {
	// set globally
	terror.SetVersion(appVersion)
	terror.SetCallbackWarn(doWarn)
	terror.SetCallbackError(doError)
	terror.SetCallbackPanic(doPanic)

	// start counting
	i := 0
	counting(&i)

	fmt.Println("finished")
}
