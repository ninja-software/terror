package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/ninja-software/terror/v2"
)

// example of simple basic error

const appVersion = "v1.2.3"

func one() (*http.Response, error) {
	// fail on purpose
	resp, err := two("http://example.commmm/")
	if err != nil {
		return nil, terror.Error(err, "get website")
	}

	return resp, nil
}

func two(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		// trace will show error msg boom
		// return nil, terror.Error(err, "boom")

		// bad! will lose tracing detail
		// return nil, terror.ErrBadContext

		// will default to use err.Error()
		// return nil, terror.New(err, "")
	}

	return resp, nil
}

func main() {
	// set globally
	terror.SetVersion(AppVersion)

	// test multi process
	for i := 0; i < 3; i++ {
		go func(num int) {
			resp, err := one()
			if err != nil {
				terror.Echo(err)
				log.Println(num, "get failed")
				return
			}
			defer resp.Body.Close()

			fmt.Println(num, "success")

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				terror.Echo(err)
				log.Println(num, "read failed")
				return
			}
			fmt.Println(string(body))
		}(i)
	}

	time.Sleep(10 * time.Second)
	fmt.Println("done")
}
