package main

import (
	"fmt"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"
)

func f(i int, jobs <-chan int, results chan<- uint64) {
	res := uint64(0)
	for j := range jobs {
		fmt.Printf("worker %d, job %d\n", i, j)

		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		res = max(res, m.Alloc)

		data := make([]byte, 1024*128)

		bigF, err := os.Open("big.mkv")
		if err != nil {
			panic(err)
		}
		_, err = bigF.Read(data)

		err = bigF.Close()

		//time.Sleep(time.Duration(4000+rand.Intn(10000)) * time.Millisecond)

		runtime.ReadMemStats(&m)
		res = max(res, m.Alloc)
		results <- res
	}
}

func human(size uint64) string {
	if size < 1024 {
		return fmt.Sprintf("%f B", float64(size))
	}
	if size < 1024*1024 {
		return fmt.Sprintf("%f KB", float64(size)/1024)
	}
	if size < 1024*1024*1024 {
		return fmt.Sprintf("%f MB", float64(size)/1024/1024)
	}
	return fmt.Sprintf("%f GB", float64(size)/1024/1024/1024)
}

func main() {
	var m runtime.MemStats

	jobs := make(chan int, 30000)
	results := make(chan uint64, 30000)

	// simulate 100 jobs
	for i := 0; i < 30000; i++ {
		jobs <- i
	}
	close(jobs)

	runtime.ReadMemStats(&m)
	fmt.Printf("allocated before jobs: %s\n", human(m.Alloc))

	now := time.Now()
	maxAllocation := uint64(0)
	// simulate 5 workers that run concurrently
	for i := 0; i < 1000; i++ {
		go f(i, jobs, results)
	}

	var allocation uint64
	for i := 0; i < 30000; i++ {
		allocation = <-results
		maxAllocation = max(maxAllocation, allocation)
	}
	elapsed := time.Since(now)
	fmt.Printf("Elapsed time: %s\n", elapsed)

	runtime.ReadMemStats(&m)
	fmt.Printf("max allocated: %s\n", human(maxAllocation))

}
