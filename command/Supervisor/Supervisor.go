//Supervisor.go
//
//Not completely functional, It only works if the command line is spoon fed
//into the function
//The functions are also not well separated and it just generally isn't what
//it should be
//Many parts of the assignment were atempted separately, but this is the only
//working model that actually compiles
package main

import (
	log "../../seelog-master/"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	//"path/filepath"
	"runtime"
	//"strconv"
	"strings"
	"sync"
)

var portFrom string
var portTo string
var portNext string

type Configuration struct {
	command	string
	output	string
	errors	string
}


func runAuthServer() {
	cmd := exec.Command("../authserver/authserver", "--log=logConfig")
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Critical(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Critical(err)
	}
	if err := cmd.Start(); err != nil {
		log.Critical(err)
	}
	b, err := json.Marshal(stderr)
	if err != nil {
		log.Errorf("error reading json: %s", err)
		os.Exit(1)
	}
	writeError := ioutil.WriteFile("out/auth.err", b, 0644)
	if writeError != nil {
		log.Errorf("error: %s", writeError)
		os.Exit(1)
	}
	o, err := json.Marshal(stdout)
	if err != nil {
		log.Errorf("error reading json: %s", err)
		os.Exit(1)
	}
	writeOut := ioutil.WriteFile("out/auth.out", o, 0644)
	if writeOut != nil {
		log.Errorf("error: %s", writeError)
		os.Exit(1)
	}
	if err := cmd.Wait(); err != nil {
		log.Critical(err)
	}
}

func runTimeServer() {
	cmd := exec.Command("../timeserver/timeserver", "--log=logConfig", "--port=8080",
              "--max-inflight=80",
              "--avg-response-ms=500", "--deviation-ms=300")
	//This isn't perfect, or even correct, but it works for now
	//next, _ := strconv.Atoi(portNext)
	//next++
	//portNext := strconv.Itoa(next)
	//_ = portNext
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Critical(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Critical(err)
	}
	if err := cmd.Start(); err != nil {
		log.Critical(err)
	}
	b, err := json.Marshal(stderr)
	if err != nil {
		log.Errorf("error reading json: %s", err)
		os.Exit(1)
	}
	writeError := ioutil.WriteFile("out/timeserver-01.err", b, 0644)
	if writeError != nil {
		log.Errorf("error: %s", writeError)
		os.Exit(1)
	}
	o, err := json.Marshal(stdout)
	if err != nil {
		log.Errorf("error reading json: %s", err)
		os.Exit(1)
	}
	writeOut := ioutil.WriteFile("out/timeserver-01.out", o, 0644)
	if writeOut != nil {
		log.Errorf("error: %s", writeError)
		os.Exit(1)
	}
	//log in case server ends unexpectedly
	if err := cmd.Wait(); err != nil {
		log.Critical(err)
	}
}


func main() {
	port := flag.String("port-range", "8080-8090", "Set the server port range, default range: 8080-8090")
	//dumpfile := flag.String("dumpfile", "backup", "This is the Supervisor dump file")
	//backupInterval := flag.Int("checkpoint-interval", 10, "This is the Supervisor backup interval")
	flag.Parse()
	ports := *port
	portRange := strings.Split(ports, "-")
	portFrom := portRange[0]
	portTo := portRange[1]
	portNext := portFrom

	fmt.Printf("port Start: %s\n", portFrom)
	fmt.Printf("port end: %s\n", portTo)
	fmt.Printf("port next: %s\n", portNext)
	//Set up for parallel proccessing
	runtime.GOMAXPROCS(2)	

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		runAuthServer()

	} ()
	
	go func() {	
		defer wg.Done()
		runTimeServer()
	} ()
	

	wg.Wait()
	/*
	
	
	
	*/
}
