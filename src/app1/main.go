package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type TaskResult struct {
	ID        int
	SleepTime int
}

func main() {
	nPtr := flag.Int("n", 0, "Количество горутин (N)")
	mPtr := flag.Int("m", 0, "Максимальное время сна в миллисекундах (M)")
	flag.Parse()

	N := *nPtr
	M := *mPtr

	if N <= 0 || M <= 0 {
		fmt.Println("Параметры N и M должны быть больше 0.")
		return
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make([]TaskResult, 0, N)

	for i := 1; i <= N; i++ {
		wg.Add(1)

		go func(goroutineID int) {
			defer wg.Done()

			sleepTime := rand.Intn(M + 1)
			time.Sleep(time.Duration(sleepTime) * time.Millisecond)

			mu.Lock()
			results = append(results, TaskResult{
				ID:        goroutineID,
				SleepTime: sleepTime,
			})
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	sort.Slice(results, func(i, j int) bool {
		return results[i].SleepTime > results[j].SleepTime
	})

	fmt.Println("Результаты работы (отсортировано по убыванию времени сна):")
	for _, res := range results {
		fmt.Printf("<%d, %d мс>\n", res.ID, res.SleepTime)
	}
}
