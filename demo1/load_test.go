package demo1

import (
	"fmt"
	"testing"
	"time"
)

// TestPutPerf load tests
func TestPutPerf(t *testing.T) {
	fmt.Printf("load testing on 1M size cache\n")
	size := 1000000
	cache := newLRUCache(size)
	start := time.Now()
	for i, td := range testData {
		if i%10000 == 0 {
			end := time.Now()
			fmt.Printf("Num: %d, Used time: %v\n", i, end.Sub(start))
			start = end
		}
		cache.Put(td.url, td.body)
	}
}
