package demo1

import (
	"fmt"
	"log"
	"testing"
	"time"
)

type testDataType struct {
	url  string
	body string
}

var testData []testDataType

const testDataSize = 10001000

// const testDataSize = 1001000

func init() {
	testData = make([]testDataType, testDataSize)
	for i := 0; i < testDataSize; i++ {
		url := fmt.Sprintf("www.fake.com/%d", i)
		testData[i] = testDataType{
			url:  url,
			body: fmt.Sprintf("This is page <b>%s</b>!", url),
		}
	}
}

func TestPutPerf(t *testing.T) {
	log.Printf("load testing... putting %d keys", testDataSize)
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
