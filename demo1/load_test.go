package demo1

import (
	"log"
	"testing"
	"time"
)

const testPerfDataSize = 1001000

func TestPutPerf(t *testing.T) {
	log.Printf("load testing... putting %d keys", testPerfDataSize)
	start := time.Now()
	size := 1000000
	cache := newLRUCache(size)
	for i, td := range testData {
		if i == size {
			log.Printf("Time to fill up cache (size %d): %v", size, time.Since(start))
		}
		cache.Put(td.url, td.body)
	}
	log.Println("Total time:", time.Since(start))
}
