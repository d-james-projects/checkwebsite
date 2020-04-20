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

func startChecking(wc Checker, website string, number int, timer time.Duration) bool {
	tz, _ := time.LoadLocation("Europe/Paris")
	start := time.Now()

	for i := 0; i < number; {
		if time.Since(start) > timer {
			parisTime := time.Now().In(tz)
			fmt.Println("Check at Paris time:", parisTime)
			stat, err := wc(website)
			if err != nil {
				fmt.Println("Website :", website, "could not be checked, possible networking issue.")
			} else {
				if stat {
					fmt.Println("Website :", website, "is Up")
				} else {
					fmt.Println("Website :", website, "is Down")
				}
			}
			start = time.Now()
			i++
		}
		time.Sleep(100 * time.Millisecond)
	}
	return true
}

func init() {
	fmt.Println("Version:\t", Version)

	flag.IntVar(&flagIterations, "i", 5, "number of times to run the webcheck")
	flag.IntVar(&flagIterations, "iterations", flagIterations, "number of times to run the webcheck")
	flag.DurationVar(&flagTimeInterval, "t", time.Duration(5*time.Second), "time (i.e. 5s) between each webcheck")
	flag.DurationVar(&flagTimeInterval, "timeinterval", flagTimeInterval, "time (i.e. 5s) between each webcheck")

	flag.Parse()
}

func main() {
	var URLs []string

	if len(flag.Args()) == 0 {
		URLs = append(URLs, "http://127.0.0.1")
	} else {
		URLs = flag.Args()
	}

	fmt.Printf("Checking the following URLs: %v\n", URLs)

	startChecking(checkWebsite, "http://127.0.0.1", flagIterations, flagTimeInterval)
}
