package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
)

// Version is used to identify the build commit for the docker image
var Version = "development"

var (
	flagIterations   int
	flagTimeInterval time.Duration
)

// Checker is a func type for controlling which check function is called
type Checker func(string) (bool, error)

func checkWebsite(website string) (bool, error) {
	isUp := false

	resp, err := http.Get(website)
	if err != nil {
		return false, err
	}

	if resp.StatusCode >= http.StatusOK &&
		resp.StatusCode < http.StatusMultipleChoices {
		isUp = true
	}

	return isUp, err
}

type result struct {
	string
	bool
	error
}

func startChecking(wc Checker, websites []string, number int, timer time.Duration) bool {
	tz, _ := time.LoadLocation("Europe/Paris")
	start := time.Now()

	rc := make(chan result)

	for i := 0; i < number; {
		if time.Since(start) > timer {
			parisTime := time.Now().In(tz)
			fmt.Println("Check at Paris time:", parisTime)
			start = time.Now()

			for _, website := range websites {
				loopWebsite := website
				go func(u string) {
					stat, err := wc(loopWebsite)
					rc <- result{loopWebsite, stat, err}
				}(website)
			}

			for c := 0; c < len(websites); c++ {
				r := <-rc

				if r.error != nil {
					fmt.Println("Website :", r.string, "could not be checked, possible networking issue.")
				} else {
					if r.bool {
						fmt.Println("Website :", r.string, "is Up")
					} else {
						fmt.Println("Website :", r.string, "is Down")
					}
				}
			}
			i++
		}
		time.Sleep(100 * time.Millisecond)
	}

	close(rc)
	return true
}

func init() {
	fmt.Println("Version:\t", Version)

	flag.IntVar(&flagIterations, "i", 5, "number of times to run the webcheck")
	flag.IntVar(&flagIterations, "iterations", flagIterations, "number of times to run the webcheck")
	flag.DurationVar(&flagTimeInterval, "t", time.Duration(5*time.Second), "time (i.e. 5s) between each webcheck")
	flag.DurationVar(&flagTimeInterval, "timeinterval", flagTimeInterval, "time (i.e. 5s) between each webcheck")
}

func main() {
	var URLs []string

	flag.Parse()

	if len(flag.Args()) == 0 {
		URLs = append(URLs, "http://127.0.0.1")
	} else {
		URLs = flag.Args()
	}

	fmt.Printf("Checking the following URLs: %v\n", URLs)

	startChecking(checkWebsite, URLs, flagIterations, flagTimeInterval)
}
