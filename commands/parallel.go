package commands

import (
	"fmt"
	"sync"
)

func RunParallelFuncs(items []string, fn func(string) error, workers uint) error {
	var wg sync.WaitGroup
	sem := make(chan bool, workers)

	for _, item := range items {
		sem <- true
		wg.Add(1)

		go func(c string) {
			defer func() { <-sem; wg.Done() }()

			err := fn(c)
			if err != nil {
				fmt.Println(" ! Failed", c, err)
			}
		}(item)
	}
	wg.Wait()
	return nil
}
