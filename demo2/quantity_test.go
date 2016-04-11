package demo2

import (
	"fmt"
	"strconv"
	"testing"
)

var testQuantities []Quantity

const (
	testDataSize     = 30000
	numOfPodsPerNode = 30
)

func init() {
	testQuantities = make([]Quantity, testDataSize)
	for i := 0; i < testDataSize; i++ {
		s := strconv.FormatInt(int64(i), 10)
		testQuantities[i] = MustParse(s)
	}
}

func ExampleQuantity() {
	q := MustParse("1Ki")
	fmt.Println(q.Value())
	q = MustParse("1")
	fmt.Println(q.MilliValue()) // 1000m
	q = MustParse("1000m")
	fmt.Println(q.Value())
	// Output:
	// 1024
	// 1000
	// 1
}

func BenchmarkMilliValue(b *testing.B) {
	for bi := 0; bi < b.N; bi++ {
		for i := 0; i < testDataSize/numOfPodsPerNode; i++ {
			c := int64(0)
			for j := 0; j < numOfPodsPerNode; j++ {
				c += testQuantities[i].MilliValue()
			}
		}
	}
}

func BenchmarkInt(b *testing.B) {
	for bi := 0; bi < b.N; bi++ {
		for i := 0; i < testDataSize/numOfPodsPerNode; i++ {
			c := 0
			for j := 0; j < numOfPodsPerNode; j++ {
				c += i*numOfPodsPerNode + j
			}
		}
	}
}
