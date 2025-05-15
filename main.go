package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	mode := os.Args[1]

	c := make(chan int)
	r := rate.Limit(1.0)
	l := rate.NewLimiter(r, 5)

	for i := range 10 {
		switch mode {
		case "allow":
			go doSomethingWithAllow(l, i, c)
		case "reserve":
			go doSomethingWithReserve(l, i, c)
		case "wait":
			go doSomethingWithWait(l, i, c)
		default:
			go doSomething(i, c)
		}
	}

	for range 10 {
		<-c
	}
}

func doSomething(x int, c chan int) {
	fmt.Printf("goroutine %d did something\n", x)

	c <- x
}

func doSomethingWithAllow(l *rate.Limiter, x int, c chan int) {
	if l.Allow() {
		fmt.Printf("Allowing %d to run\n", x)
	}

	c <- x
}

func doSomethingWithWait(l *rate.Limiter, x int, c chan int) {
	err := l.Wait(context.Background())
	if err != nil {
		fmt.Printf("Error waiting for %d: %v\n", x, err)
		c <- x
		return
	}

	fmt.Printf("Allowing %d to run\n", x)
	c <- x
}

func doSomethingWithReserve(l *rate.Limiter, x int, c chan int) {
	r := l.Reserve()
	if !r.OK() {
		return
	}

	fmt.Printf("Reserving %d to run\n", x)
	d := r.Delay()
	time.Sleep(d)
	fmt.Printf("Allowing %d to run\n", x)

	c <- x
}

func doSomethingWithCircuitBreaker() {
	s := rate.Sometimes{Every: 2}
	s.Do(func() { fmt.Println("1") })
	s.Do(func() { fmt.Println("2") })
	s.Do(func() { fmt.Println("3") })
}
