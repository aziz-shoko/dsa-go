package main

import (
	"io"
	"os"
	"fmt"
	"time"
)

// Take a thin slice of functionality and make it work end-to-end, backed by tests.
// So first make tests pass, and make it work end-to-end by actually making it work with 
// go run command or something

const finalWord = "Go!"
const coundownStart = 3

type Sleeper interface {
	Sleep()
}

type ConfigurableSleeper struct {
	duration time.Duration
	sleep    func(time.Duration)
}

func (c *ConfigurableSleeper) Sleep() {
	c.sleep(c.duration)
}

func main() {
	sleeper := &ConfigurableSleeper{1 * time.Second, time.Sleep}
	Countdown(os.Stdout, sleeper)
}

func Countdown(out io.Writer, sleep Sleeper) {
	for i := coundownStart; i > 0; i-- {
		fmt.Fprintln(out, i)
		sleep.Sleep()
	}	
	fmt.Fprint(out, finalWord)
}