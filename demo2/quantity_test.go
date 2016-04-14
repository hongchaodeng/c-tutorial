package demo2

import (
	"fmt"
	"strconv"
	"testing"

	"speter.net/go/exp/math/dec/inf"
)

var _ = inf.Dec{}

var testQuantities []Quantity

const (
	numOfPodsPerNode = 30
	numOfNodes       = 1000
	testDataSize     = numOfNodes * numOfPodsPerNode
)

func init() {
	testQuantities = make([]Quantity, testDataSize)
	for i := 0; i < testDataSize; i++ {
		s := strconv.FormatInt(int64(i), 10)
		testQuantities[i] = MustParse(s)
	}
}

func ExampleQuantity() {
	// This is to show how we parse value from user config into "Quantity" API.
	// Quantity provides the ability to transform into different units.
	// So we can get value in basic unit, milli unit, or others.
	q := MustParse("1Ki")
	fmt.Println(q.Value())
	q = MustParse("1")
	fmt.Println(q.MilliValue())
	q = MustParse("1000m")
	fmt.Println(q.Value())
	// Output:
	// 1024
	// 1000
	// 1
}

func BenchmarkMilliValue(b *testing.B) {
	// This is to simulate the same problem we have previously in aggregating values.
	// Basically we have serious latency here like 30ms.
	// After optimization, it's <1ms.
	// We use it to reduce scheduling latency from 57ms to 25ms, doubling scheduling rate.
	for bi := 0; bi < b.N; bi++ {
		for i := 0; i < numOfNodes; i++ {
			// aggregate resources per node
			c := int64(0)
			for j := 0; j < numOfPodsPerNode; j++ {
				c += testQuantities[i].MilliValue()
			}
		}
	}
}

// func BenchmarkInt(b *testing.B) {
// 	for bi := 0; bi < b.N; bi++ {
// 		for i := 0; i < numOfNodes; i++ {
// 			c := 0
// 			for j := 0; j < numOfPodsPerNode; j++ {
// 				c += i*numOfPodsPerNode + j
// 			}
// 		}
// 	}
// }

// func BenchmarkLargeValuePool(b *testing.B) {
// 	s := big.NewInt(math.MaxInt64)
// 	s.Mul(s, big.NewInt(1000))
// 	for bi := 0; bi < b.N; bi++ {
// 		scaledValue(s, 10, 3)
// 	}
// }

// func BenchmarkLargeValueInf(b *testing.B) {
// 	tmp := &inf.Dec{}
// 	q := MustParse("1")
// 	for bi := 0; bi < b.N; bi++ {
// 		tmp.Round(tmp.Mul(q.Amount, decThousand), 0, inf.RoundUp).UnscaledBig().Int64()
// 	}
// }
