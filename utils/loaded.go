package utils

import (
	"fmt"
	"time"

	"github.com/paulbellamy/ratecounter"
)

type Loaded struct {
	rate *ratecounter.RateCounter
}

func NewLoaded() *Loaded {
	return &Loaded{
		rate: ratecounter.NewRateCounter(1 * time.Second),
	}
}

func (c *Loaded) Check() {
	c.rate.Incr(1)
}

func (c *Loaded) StartLogging() {
	go func() {
		for {
			fmt.Printf("per/s: %v\n", c.rate.Rate())
			time.Sleep(time.Second)
		}
	}()
}
