package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestRandomValue(t *testing.T) {
	numTests := 100000
	minVal := 1
	maxVal := 100
	expectedNumberOfHits := numTests / ((maxVal - minVal) + 1)
	criticalChi := 123.225 // This is for a 5% with a DF of 99

	totalValues := make(map[int]int)

	for i := 0; i < numTests; i++ {
		val := randomValue(minVal, maxVal)
		totalValues[val] = totalValues[val] + 1
	}

	keys := make([]int, 0, len(totalValues))
	for k := range totalValues {
		keys = append(keys, k)
	}

	for j := minVal; j < maxVal; j++ {
		_, ok := totalValues[j]
		if !ok {
			t.Errorf("%d value was not reached", j)
		}
	}

	sort.Ints(keys)

	chiSquared := 0.0
	for _, val := range keys {
		numHits := totalValues[val]
		subChi := math.Pow((float64(numHits)-float64(expectedNumberOfHits)), 2) / float64(expectedNumberOfHits)
		chiSquared = chiSquared + subChi
		fmt.Printf("%d:\n", val)
		fmt.Printf("    numHits  = %d\n", numHits)
		fmt.Printf("    expected = %d\n", expectedNumberOfHits)
		fmt.Printf("    subChi   = %.6f\n", subChi)
	}

	fmt.Printf("\nchi-squared = %.6f\n", chiSquared)
	fmt.Printf("critical-ch = %.6f\n", criticalChi)
	if chiSquared > criticalChi {
		t.Errorf("the chi-squared value was too large")
	}

}
