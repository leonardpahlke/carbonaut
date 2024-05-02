package carbonautserver

import "time"

func Listen(exitChan chan int) {
	// TODO: start server and listen on endpoints
	// TODO: add some kind of http server which can be used to stop carbonaut, update the configuration, get a status etc.
	time.Sleep(100 * time.Second)
	exitChan <- 1
}
