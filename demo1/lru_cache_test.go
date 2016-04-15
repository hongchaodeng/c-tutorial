package demo1

import (
	"fmt"
	"testing"
)

type testDataType struct {
	url  string
	body string
}

var testData []testDataType

const testDataSize = 10000100

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
func TestBasics(t *testing.T) {
	cache := newLRUCache(2)
	cache.Put("1", "1")
	cache.Put("2", "2")
	res, ok := cache.Get("1")
	if !ok || res != "1" {
		t.Errorf("unexpected!")
	}
	res, ok = cache.Get("2")
	if !ok || res != "2" {
		t.Errorf("unexpected!")
	}

	cache.Put("3", "3")
	if _, ok := cache.Get("1"); ok {
		t.Errorf("unexpected!")
	}
	cache.Put("4", "4")
	if _, ok := cache.Get("2"); ok {
		t.Errorf("unexpected!")
	}
}

// ---------------------------------
// 问一下大家知不知道 go tool pprof
// ---------------------------------

func BenchmarkCacheFull1KSize(b *testing.B) {
	benchmarkCacheFull(b, 1000)
}

func BenchmarkCacheFull1MSize(b *testing.B) {
	benchmarkCacheFull(b, 1000000)
}
func BenchmarkCacheFull10MSize(b *testing.B) {
	benchmarkCacheFull(b, 10000000)
}

// benchmarkCacheFull tests how long it takes to put a new key
// when LRU cache is full
func benchmarkCacheFull(b *testing.B, size int) {
	cache := newLRUCache(size)
	for i := 0; i < size; i++ {
		cache.Put(testData[i].url, testData[i].body)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := (i + size) % testDataSize
		cache.Put(testData[c].url, testData[c].body)
	}
}
