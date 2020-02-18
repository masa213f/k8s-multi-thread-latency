package main

import (
	"crypto/sha512"
	"flag"
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	procs := flag.Int("procs", runtime.NumCPU(), "GOMAXPROCS")
	sleep := flag.Duration("sleep", time.Millisecond, "sleep between iterations")
	iterations := flag.Int("iterations", 100, "number of iterations")
	flag.Parse()

	fmt.Printf("procs: %d, sleep: %s, iterations: %d\n", *procs, (*sleep).String(), *iterations)
	runtime.GOMAXPROCS(*procs)

	type result struct {
		ms    int
		start time.Time
		end   time.Time
	}
	resultList := make([]*result, *iterations, *iterations)

	start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < *iterations; i++ {
		time.Sleep(*sleep)

		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s := time.Now()
			burn(time.Millisecond * 5)
			e := time.Now()
			t := ms(e.Sub(s))
			resultList[i] = &result{ms: t, start: s, end: e}
		}(i)
	}
	wg.Wait()
	end := time.Now()

	fmt.Printf("total, %5dms, %s, %s\n", ms(end.Sub(start)), timeToString(start), timeToString(end))
	for i := 0; i < *iterations; i++ {
		fmt.Printf("%5d, %5dms, %s, %s\n", i, resultList[i].ms, timeToString(resultList[i].start), timeToString(resultList[i].end))
	}
}

func ms(duration time.Duration) int {
	return int(duration.Nanoseconds() / 1000 / 1000)
}

func timeToString(t time.Time) string {
	// Unify the number of displayed characters as follows.
	// "2006-01-02T15:04:05.999"
	// Because Time.Format() sometimes fluctuates the number of displayed characters.
	nano := t.Format(".999")
	if len(nano) == 0 {
		nano = ".000"
	} else {
		nano = (nano + "000")[:4]
	}
	return fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d%s",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), nano)
}

func burn(duration time.Duration) {
	s := time.Now()

	for {
		sum := sha512.New()
		sum.Write([]byte("banana"))
		sum.Sum([]byte{})

		if time.Since(s) > duration {
			break
		}
	}
}
