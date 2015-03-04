// Thank you Professor Bernstein!

package main

import (

	"fmt"
	"net/http"
	"time"
	"github.com/JKowalsky/counter"
	log "github.com/cihub/seelog"
)

var (
	rate = 200
	burst = 30
	timeoutMs = 400
	runtime = 20 * time.Second
	url = "localhost:8080/time"
)
var (
	c = counter.New()

)

var convert = map[int]string {
	1: "100s",
	2: "200s",
	3: "300s",
	4: "400s",
	5: "500s",
}


func request() {
	log.Info("New Request.")
	timeout := time.Duration(timeoutMs) * time.Millisecond
	client := http.Client{
		Timeout : timeout,
	}
	response, err := client.Get(url)
	if err != nil {
		log.Error("No response.")
		c.Incr("total", 1)
		return
	}
	key, ok := convert[response.StatusCode / 100]
	log.Info("Response: ", key)
	if !ok {
		log.Error("Response was an error.")
		key = "errors"
	}
	c.Incr(key, 1)
}

func load() {
	timeout := time.Tick(runtime)
	interval := time.Duration((1000000 * burst) / rate) * time.Microsecond
	period := time.Tick(interval)

	log.Info("Loading.")
	for {
		log.Info("Fire a burst.")
		// fire off burst
		for i := 0; i < burst; i++ {
			go request()
		}
		//wait for next tick
		<- period

		// poll for timeout
		select {
		case <-timeout :
			return
		default:

		}
		
	}
}

func main () {
	load()
	time.Sleep( time.Duration( 2 * timeoutMs) * time.Millisecond )
	fmt.Printf("total: \t%d\n", c.Get("total"))
	fmt.Printf("100s total: \t%d\n", c.Get("100s"))
	fmt.Printf("200s total: \t%d\n", c.Get("200s"))
	fmt.Printf("300s total: \t%d\n", c.Get("300s"))
	fmt.Printf("400s total: \t%d\n", c.Get("400s"))
	fmt.Printf("500s total: \t%d\n", c.Get("500s"))
	fmt.Printf("total: \t%d\n", c.Get("total"))


}
