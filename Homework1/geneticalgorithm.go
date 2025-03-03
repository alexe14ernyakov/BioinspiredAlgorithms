package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const (
	tournamentSize  = 3
	popSize         = 200
	crossProb       = 0.7
	mutProb         = 0.1
	xMin            = 2.0
	xMax            = 4.0
	stagnationLimit = 20
)

func randChoice[T any](arr []T) T {
	return arr[rand.Intn(len(arr))]
}

func fitness(x float64) float64 {
	numerator := math.Cos(math.Exp(x))
	denominator := math.Sin(math.Log(x))
	return numerator / denominator
}

func genIndividual() float64 {
	return xMin + (xMax-xMin)*rand.Float64()
}

func genPopulation() (res []float64) {
	res = make([]float64, popSize)

	for i := range res {
		res[i] = genIndividual()
	}
	return
}

func tournamentSelection(population []float64) (best float64) {
	best = randChoice(population)
	for i := 1; i < tournamentSize; i++ {
		contender := randChoice(population)
		if fitness(contender) > fitness(best) {
			best = contender
		}
	}
	return
}

// Одноточечный кроссинговер
func crossingover(p1, p2 float64) float64 {
	if rand.Float64() < crossProb {
		return (p1 + p2) / 2
	}
	return p1
}

func mutate(ind float64) float64 {
	if rand.Float64() < mutProb {
		delta := (rand.Float64() - 0.5) * 0.1
		mutant := ind + delta
		if mutant < xMin {
			mutant = xMin
		} else if mutant > xMax {
			mutant = xMax
		}
		return mutant
	}
	return ind
}

func main() {
	for range 10 {
		rand.Seed(time.Now().UnixNano())
		startTime := time.Now()

		population := genPopulation()
		bestIndividual := population[0]
		maxExtremum := fitness(bestIndividual)
		stagnationCount := 0
		generation := 0

		for stagnationCount < stagnationLimit {
			newPopulation := genPopulation()

			for i := range newPopulation {
				p1 := tournamentSelection(population)
				p2 := tournamentSelection(population)
				child := crossingover(p1, p2)
				child = mutate(child)
				newPopulation[i] = child
			}
			population = newPopulation

			improved := false
			for _, ind := range population {
				value := fitness(ind)
				if value > maxExtremum {
					bestIndividual = ind
					maxExtremum = value
					improved = true
				}
			}

			if improved {
				stagnationCount = 0
			} else {
				stagnationCount++
			}

			fmt.Printf("Поколение %d: x = %.10f; f(x) = %.10f\n", generation, bestIndividual, maxExtremum)
			generation++
		}

		workTime := time.Since(startTime)
		fmt.Printf("Лучшее решение, найденное алгоритмом: f(%.10f) = %.10f\n", bestIndividual, maxExtremum)
		fmt.Printf("Затраченное время: %d мс\n", workTime.Milliseconds())
	}
}
