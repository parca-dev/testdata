package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: ./a.out runtime_in_ms")
		return
	}

	dur, err := strconv.Atoi(os.Args[1])
	fmt.Println(dur)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	runFor := time.Duration(dur) * time.Millisecond

	fmt.Println("Running with pid", os.Getpid())
	futureTime := time.Now().Add(runFor)
	for {
		if time.Now().After(futureTime) {
			break
		}
	}
	fmt.Println("Exiting")
}
