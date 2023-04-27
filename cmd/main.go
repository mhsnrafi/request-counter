package main

import (
	"log"
	"time"

	"github.com/mhsnrafi/request-counter/counter"
	"github.com/mhsnrafi/request-counter/server"
)

func main() {
	config := &counter.Config{
		TimeWindow:      60 * time.Second,
		PersistenceFile: "counts.json",
	}

	cnt := counter.NewCounter(config)
	err := cnt.Load()
	if err != nil {
		log.Fatal(err)
	}

	server.Run(cnt)
}
