// monitor.go

package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"flag"
	"io/ioutil"
	"strings"
	//"os"

	//"sync"
	"time"
)

// list of URLs to monitor
var targets = "http://localhost:8080/monitor, http://localhost:9090/monitor"

// sliced URL from targets
var server1 = "http://localhost:8080/monitor"
var server2 = "http://localhost:9090/monitor"

// specify the interval (in seconds) between sample requests.
var sample_interval_sec = 1.75

// run time of monitor program.
var runtime_sec = 5.0

//These structs and maps are here in case I figure out how to get everything working
type Sample struct {
	time	int
	value	int
}

type Counter struct {
	dictionary	[]byte
}

//var sequence map[string][]Counter

//var monitors map[string][]sequence

func monServer(server string) {

	runTime := time.Duration(runtime_sec) * time.Second
	period := time.Tick(runTime)

	// Get monitor data from server
	for {
		res, err := http.Get(server)
    		if err != nil {
        		panic(err.Error())
    		}

    		body, err := ioutil.ReadAll(res.Body)
    		if err != nil {
        		panic(err.Error())
    		}
		
		fmt.Println("Server: ", server)

		g := CToGoString(body[:])
		
		//This is here as a test to check that data has been input
		fmt.Println(g)
		fmt.Println("\n")

		<-period
	}
}

//http://stackoverflow.com/questions/14230145/what-is-the-best-way-to-convert-byte-array-to-string
// convert a byte array to a string
func CToGoString(c []byte) string {
	n := -1
	for i, b := range c {
		if b == 0 {
			break
		}
		n = i
	}
	return string(c[:n+1])
}

func sliceURL(urls string) {
	sliced := strings.Split(urls, ", ")

	if len(sliced) != 2 {
		return
	}

	// update server names
	server1 = sliced[0]
	server2 = sliced[1]

}

func main() {
	flag.StringVar(&targets, "targets", "two URLs", "The two URLs that are monitored")
	flag.Float64Var(&sample_interval_sec, "sample-interval-sec", 1.75, "sampling interval in sec")
	flag.Float64Var(&runtime_sec, "runtime-sec", 1.75, "runtime of Monitor in sec")

	flag.Parse()

	sliceURLstring(targets)

	go monServer(server1)
	go monServer(server2)

	sleepTime := time.Duration(int(sample_interval_sec)) * time.Second
	time.Sleep(sleepTime)

}
