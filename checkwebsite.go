package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Version is used to identify the build commit for the docker image
var Version = "development"

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

func main() {
	fmt.Println("Version:\t", Version)

	checkURL := "https://www.sky.com/"
	number := 5
	timer := time.Duration(5 * time.Second)

	// app cmd line has 3 args if none are supplied the defaults apply (see above)
	// checkwebsite <url> <times to check/int> <delay between checks/int secs>
	if len(os.Args) > 4 {
		fmt.Println("Wrong number of parameters supplied, requires just one website url for the check.")
		os.Exit(1)
	}

	if len(os.Args) == 2 {
		checkURL = os.Args[1]
	} else if len(os.Args) == 3 {
		checkURL = os.Args[1]
		i, _ := strconv.Atoi(os.Args[2])
		number = i
	} else if len(os.Args) == 4 {
		checkURL = os.Args[1]
		i, _ := strconv.Atoi(os.Args[2])
		number = i
		t, _ := strconv.Atoi(os.Args[3])
		timer = time.Duration(time.Duration(t) * time.Second)
	}

	startChecking(checkWebsite, checkURL, number, timer)
}
