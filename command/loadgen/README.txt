CSS490 Assignment 5
By Joel Richart


loadgen.go contains following flags 

	- rate        : Set average rate of requests (per second) 
	- burst       : Set number of concurrent requests to issue )
	- timeout-ms  : Set max time to wait for response ( in microseconds )
	- runtime     : Set number of seconds to process
	- url         : set URL to Sample

if flags are not provided in the command-line, default settings are set to

rete: 200
burst: 20
timeout-ms : 1000
runtime : 10
url : http://localhost:8080/time




