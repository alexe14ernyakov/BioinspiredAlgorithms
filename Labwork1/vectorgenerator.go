package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const (
	vectorLength = 24
	numVectors   = 50
)

func generateKnapsackVector(maxValue int) []int {
	vector := make([]int, vectorLength)
	for i := 0; i < vectorLength; i++ {
		vector[i] = rand.Intn(maxValue) + 1
	}
	return vector
}

func vectorKey(vec []int) string {
	key := ""
	for _, v := range vec {
		key += fmt.Sprintf("%d,", v)
	}
	return key
}

func generateUniqueVectors() [][]int {
	maxValue := int(math.Pow(2, 24/1.4))
	vectorsMap := make(map[string]bool)
	vectors := [][]int{}

	for len(vectors) < numVectors {
		vec := generateKnapsackVector(maxValue)
		key := vectorKey(vec)
		if !vectorsMap[key] {
			vectorsMap[key] = true
			vectors = append(vectors, vec)
		}
	}
	return vectors
}

func main() {
	rand.Seed(time.Now().UnixNano())

	vectors := generateUniqueVectors()

	for i := 0; i < len(vectors); i++ {
		fmt.Printf("%d: %v\n", i+1, vectors[i])
	}
}
