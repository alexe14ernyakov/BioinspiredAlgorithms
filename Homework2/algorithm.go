package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const (
	dim        = 2     // Размерность задачи
	numFish    = 30    // Размер популяции
	iterations = 100   // Количество итераций
	stepInd    = 0.1   // Индивидуальный шаг
	stepVol    = 0.01  // Волитивное движение
	boundMin   = -5.12 // Минимум области поиска
	boundMax   = 5.12  // Максимум области поиска
)

type Fish struct {
	position      []float64
	fitness       float64
	mass          float64
	deltaPosition []float64
}

func rastrigin(x []float64) float64 {
	A := 10.0
	sum := A * float64(len(x))
	for _, xi := range x {
		sum += xi*xi - A*math.Cos(2*math.Pi*xi)
	}
	return sum
}

func randomVector(n int, min, max float64) []float64 {
	v := make([]float64, n)
	for i := range v {
		v[i] = min + rand.Float64()*(max-min)
	}
	return v
}

func clampVector(v []float64, min, max float64) {
	for i := range v {
		if v[i] < min {
			v[i] = min
		}
		if v[i] > max {
			v[i] = max
		}
	}
}

func fishSchoolSearch() ([]float64, float64) {
	rand.Seed(time.Now().UnixNano())
	school := make([]Fish, numFish)
	var bestPosition []float64
	bestFitness := math.MaxFloat64

	for i := range school {
		pos := randomVector(dim, boundMin, boundMax)
		fit := rastrigin(pos)
		school[i] = Fish{
			position:      pos,
			fitness:       fit,
			mass:          1.0,
			deltaPosition: make([]float64, dim),
		}
		if fit < bestFitness {
			bestFitness = fit
			bestPosition = append([]float64{}, pos...)
		}
	}

	for iter := 0; iter < iterations; iter++ {
		totalWeightGain := 0.0

		for i := range school {
			direction := randomVector(dim, -1, 1)
			newPos := make([]float64, dim)
			for j := range newPos {
				newPos[j] = school[i].position[j] + direction[j]*stepInd
			}
			clampVector(newPos, boundMin, boundMax)

			newFit := rastrigin(newPos)
			if newFit < school[i].fitness {
				for j := range school[i].position {
					school[i].deltaPosition[j] = newPos[j] - school[i].position[j]
					school[i].position[j] = newPos[j]
				}
				weightGain := school[i].fitness - newFit
				school[i].fitness = newFit
				school[i].mass += weightGain
				totalWeightGain += weightGain
			} else {
				for j := range school[i].deltaPosition {
					school[i].deltaPosition[j] = 0
				}
			}
		}

		totalMass := 0.0
		for i := range school {
			if school[i].mass < 1.0 {
				school[i].mass = 1.0
			}
			if school[i].mass > 5.0 {
				school[i].mass = 5.0
			}
			totalMass += school[i].mass
		}

		collectiveMove := make([]float64, dim)
		for i := range school {
			for j := range collectiveMove {
				collectiveMove[j] += school[i].deltaPosition[j] * school[i].mass
			}
		}
		for j := range collectiveMove {
			collectiveMove[j] /= totalMass
		}
		for i := range school {
			for j := range school[i].position {
				school[i].position[j] += collectiveMove[j]
			}
			clampVector(school[i].position, boundMin, boundMax)
			school[i].fitness = rastrigin(school[i].position)
		}

		barycenter := make([]float64, dim)
		for i := range school {
			for j := range barycenter {
				barycenter[j] += school[i].position[j] * school[i].mass
			}
		}
		for j := range barycenter {
			barycenter[j] /= totalMass
		}
		for i := range school {
			for j := range school[i].position {
				diff := school[i].position[j] - barycenter[j]
				if totalWeightGain > 0 {
					school[i].position[j] -= stepVol * rand.Float64() * diff
				} else {
					school[i].position[j] += stepVol * rand.Float64() * diff
				}
			}
			clampVector(school[i].position, boundMin, boundMax)
			school[i].fitness = rastrigin(school[i].position)
		}

		for _, f := range school {
			if f.fitness < bestFitness {
				bestFitness = f.fitness
				bestPosition = append([]float64{}, f.position...)
			}
		}

		fmt.Printf("%3d | Best fitness: %.6f\n", iter+1, bestFitness)
	}

	return bestPosition, bestFitness
}

func main() {
	bestPos, bestVal := fishSchoolSearch()
	fmt.Println("\nBest position:", bestPos)
	fmt.Println("Function value:", bestVal)
}
