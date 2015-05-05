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
	"strconv"
	"strings"
	"sync"
)

var portFrom int
var portTo int
var portNext int

type Configuration struct {
	command	string
	output	string
	errors	string
}

func runAuthServer(authConfig Configuration) {
	cmd := exec.Command(authConfig.command)
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
	writeError := ioutil.WriteFile(authConfig.errors, b, 0644)
	if writeError != nil {
		log.Errorf("error: %s", writeError)
		os.Exit(1)
	}
	o, err := json.Marshal(stdout)
	if err != nil {
		log.Errorf("error reading json: %s", err)
		os.Exit(1)
	}
	writeOut := ioutil.WriteFile(authConfig.output, o, 0644)
	if writeOut != nil {
		log.Errorf("error: %s", writeError)
		os.Exit(1)
	}
	if err := cmd.Wait(); err != nil {
		log.Critical(err)
	}
}

func runTimeServer(timeConfig Configuration) {
	cmd := exec.Command(timeConfig.command)
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
	writeError := ioutil.WriteFile(timeConfig.errors, b, 0644)
	if writeError != nil {
		log.Errorf("error: %s", writeError)
		os.Exit(1)
	}
	o, err := json.Marshal(stdout)
	if err != nil {
		log.Errorf("error reading json: %s", err)
		os.Exit(1)
	}
	writeOut := ioutil.WriteFile(timeConfig.output, o, 0644)
	if writeOut != nil {
		log.Errorf("error: %s", writeError)
		os.Exit(1)
	}
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
	portStart := portRange[0]
	portEnd := portRange[1]
	portFrom, _ := strconv.Atoi(portStart)
	portTo, _ := strconv.Atoi(portEnd)
	portNext := portFrom

	fmt.Printf("From Port is: \t%d\n", portFrom)
	fmt.Printf("To Port is: \t%d\n", portTo)
	fmt.Printf("Next Port is: \t%d\n", portNext)

	//./supervisor "$(<config.json)"
	filename := os.Args[1]
	
	file, fileErr := ioutil.ReadFile(filename)
	if fileErr != nil {
		log.Errorf("file error: %s", fileErr)
		return
	}

	dec := json.NewDecoder(strings.NewReader(file))
	var c Configuration

	if err := dec.Decode(&c); err == io.EOF {
		break
	} else if err != nil {
		log.Critical(err)
	}
	authServerConfig := c
	//Set up for parallel proccessing
	runtime.GOMAXPROCS(3)	

	var wg sync.WaitGroup
	wg.Add(3)

	go func(authServerConfig Configuration) {
		defer wg.Done()
		runAuthServer(authServerConfig)

	} ()
	if err := dec.Decode(&c); err == io.EOF {
		break
	} else if err != nil {
		log.Critical(err)
	}
	timeServerConfig1 := c
	go func(timeServerConfig Configuration) {	
		defer wg.Done()
		runTimeServer(timeServerConfig1)
	} ()
	if err := dec.Decode(&c); err == io.EOF {
		break
	} else if err != nil {
		log.Critical(err)
	}
	timeServerConfig2 := c
	go func(timeServerConfig Configuration) {	
		defer wg.Done()
		runTimeServer(timeServerConfig2)
	} ()

	wg.Wait()
	/*
	
	
	
	*/
}
