package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)

type Item struct {
	Weight int
	Index  int
}

type KnapsackProblem struct {
	ID          int
	Target      int
	Ratio       float64
	BruteTimeMs float64
}

type Chromosome struct {
	Genes   []bool
	Fitness int
	Weight  int
}

type GAConfig struct {
	PopulationSize   int
	MutationRate     float64
	CrossoverRate    float64
	MaxGenerations   int
	MaxNoImprovement int
}

type GAResult struct {
	VectorID          int
	ProblemID         int
	TargetWeight      int
	AchievedWeight    int
	Fitness           int
	Generations       int
	DurationMs        float64
	TerminationReason string
	BestSolution      []int
}

func main() {
	rand.Seed(time.Now().UnixNano())

	itemsList, err := readItems("knapsack_vectors.csv")
	if err != nil {
		log.Fatal("Error reading items:", err)
	}

	problemsList, err := readProblems("problems.csv")
	if err != nil {
		log.Fatal("Error reading problems:", err)
	}

	bruteTimes, err := readBruteTimes("knapsack_solutions_ms.csv")
	if err != nil {
		log.Println("Warning: could not read brute times, using default values:", err)
		bruteTimes = make([][]float64, len(itemsList))
		for i := range bruteTimes {
			bruteTimes[i] = make([]float64, 15)
		}
	}

	config := GAConfig{
		PopulationSize:   1000,
		MutationRate:     0.05,
		CrossoverRate:    0.7,
		MaxGenerations:   100,
		MaxNoImprovement: 2,
	}

	resultsFile, err := os.Create("ga_solutions.csv")
	if err != nil {
		log.Fatal("Error creating results file:", err)
	}
	defer resultsFile.Close()

	writer := csv.NewWriter(resultsFile)
	defer writer.Flush()

	header := []string{
		"VectorID", "ProblemID", "TargetWeight", "AchievedWeight",
		"Fitness", "Generations", "DurationMs", "TerminationReason", "SolutionItems",
	}
	writer.Write(header)

	for vectorID, items := range itemsList {
		problems := problemsList[vectorID]
		for problemID, problem := range problems {
			problem.BruteTimeMs = bruteTimes[vectorID][problemID]

			startTime := time.Now()
			bestSolution, stats := geneticAlgorithm(items, problem, config)
			duration := time.Since(startTime).Seconds() * 1000

			result := GAResult{
				VectorID:          vectorID + 1,
				ProblemID:         problem.ID,
				TargetWeight:      problem.Target,
				AchievedWeight:    bestSolution.Weight,
				Fitness:           bestSolution.Fitness,
				Generations:       stats["generations"].(int),
				DurationMs:        duration,
				TerminationReason: stats["termination_reason"].(string),
				BestSolution:      getSolutionIndices(bestSolution, items),
			}

			record := []string{
				strconv.Itoa(result.VectorID),
				strconv.Itoa(result.ProblemID),
				strconv.Itoa(result.TargetWeight),
				strconv.Itoa(result.AchievedWeight),
				strconv.Itoa(result.Fitness),
				strconv.Itoa(result.Generations),
				fmt.Sprintf("%.3f", result.DurationMs),
				result.TerminationReason,
				formatSolution(result.BestSolution),
			}
			writer.Write(record)

			fmt.Printf("Вектор %d Задача %d: , Достигнуто: %d (Целевой вес: %d), Фитнесс-функция = %d, Поколений: %d, Время работы: %.2fms, Причина остановки: %s\n",
				result.VectorID, result.ProblemID, result.AchievedWeight, result.TargetWeight,
				result.Fitness, result.Generations, result.DurationMs, result.TerminationReason)
		}
	}
}

func geneticAlgorithm(items []Item, problem KnapsackProblem, config GAConfig) (Chromosome, map[string]interface{}) {
	startTime := time.Now()
	stats := map[string]interface{}{
		"generations":        0,
		"termination_reason": "max_generations",
	}

	population := initializePopulation(len(items), config.PopulationSize, items, problem.Target)
	bestSolution := findBest(population)
	noImprovementCount := 0

	for gen := 0; gen < config.MaxGenerations; gen++ {
		stats["generations"] = gen + 1

		newPopulation := make([]Chromosome, 0, config.PopulationSize)
		for len(newPopulation) < config.PopulationSize {
			parent1 := tournamentSelection(population, 3)
			parent2 := tournamentSelection(population, 3)

			var child1, child2 Chromosome
			if rand.Float64() < config.CrossoverRate {
				child1, child2 = crossover(parent1, parent2)
			} else {
				child1, child2 = parent1, parent2
			}

			child1 = mutate(child1, config.MutationRate)
			child2 = mutate(child2, config.MutationRate)
			child1 = calculateFitness(child1, items, problem.Target)
			child2 = calculateFitness(child2, items, problem.Target)
			newPopulation = append(newPopulation, child1, child2)
		}

		population = newPopulation
		currentBest := findBest(population)

		if currentBest.Fitness == 0 {
			stats["termination_reason"] = "zero_fitness"
			return currentBest, stats
		}

		if currentBest.Fitness >= bestSolution.Fitness {
			noImprovementCount++
			if noImprovementCount >= config.MaxNoImprovement {
				stats["termination_reason"] = "no_improvement"
				return bestSolution, stats
			}
		} else {
			bestSolution = currentBest
			noImprovementCount = 0
		}

		elapsed := time.Since(startTime).Seconds() * 1000
		if problem.BruteTimeMs > 0 && elapsed >= 2*problem.BruteTimeMs {
			stats["termination_reason"] = "time_exceeded"
			return bestSolution, stats
		}
	}

	return bestSolution, stats
}

func calculateFitness(c Chromosome, items []Item, target int) Chromosome {
	totalWeight := 0
	for i, gene := range c.Genes {
		if gene {
			totalWeight += items[i].Weight
		}
	}

	c.Weight = totalWeight
	c.Fitness = int(math.Abs(float64(target - totalWeight)))
	return c
}

func initializePopulation(itemCount, populationSize int, items []Item, target int) []Chromosome {
	population := make([]Chromosome, populationSize)
	for i := range population {
		genes := make([]bool, itemCount)
		for j := range genes {
			genes[j] = rand.Float32() < 0.35
		}
		population[i] = calculateFitness(Chromosome{Genes: genes}, items, target)
	}
	return population
}

func tournamentSelection(population []Chromosome, tournamentSize int) Chromosome {
	best := population[rand.Intn(len(population))]
	for i := 1; i < tournamentSize; i++ {
		contender := population[rand.Intn(len(population))]
		if contender.Fitness < best.Fitness {
			best = contender
		}
	}
	return best
}

func crossover(parent1, parent2 Chromosome) (Chromosome, Chromosome) {
	if len(parent1.Genes) != len(parent2.Genes) {
		return parent1, parent2
	}

	crossoverPoint := rand.Intn(len(parent1.Genes))
	child1Genes := make([]bool, len(parent1.Genes))
	child2Genes := make([]bool, len(parent2.Genes))

	for i := 0; i < crossoverPoint; i++ {
		child1Genes[i] = parent1.Genes[i]
		child2Genes[i] = parent2.Genes[i]
	}
	for i := crossoverPoint; i < len(parent1.Genes); i++ {
		child1Genes[i] = parent2.Genes[i]
		child2Genes[i] = parent1.Genes[i]
	}

	return Chromosome{Genes: child1Genes}, Chromosome{Genes: child2Genes}
}

func mutate(c Chromosome, mutationRate float64) Chromosome {
	for i := range c.Genes {
		if rand.Float64() < mutationRate {
			c.Genes[i] = !c.Genes[i]
		}
	}
	return c
}

func findBest(population []Chromosome) Chromosome {
	best := population[0]
	for _, c := range population {
		if c.Fitness < best.Fitness {
			best = c
		}
	}
	return best
}

func getSolutionIndices(c Chromosome, items []Item) []int {
	var indices []int
	for i, gene := range c.Genes {
		if gene {
			indices = append(indices, items[i].Index)
		}
	}
	sort.Ints(indices)
	return indices
}

func formatSolution(indices []int) string {
	if len(indices) == 0 {
		return ""
	}
	str := fmt.Sprintf("%v", indices)
	return str[1 : len(str)-1]
}

func readItems(filename string) ([][]Item, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	itemsList := make([][]Item, len(records))
	for i, record := range records {
		items := make([]Item, len(record))
		for j, value := range record {
			weight, err := strconv.Atoi(value)
			if err != nil {
				return nil, err
			}
			items[j] = Item{Weight: weight, Index: j}
		}
		itemsList[i] = items
	}

	return itemsList, nil
}

func readProblems(filename string) ([][]KnapsackProblem, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	problemsList := make([][]KnapsackProblem, 0)
	currentGroup := make([]KnapsackProblem, 0, 15)

	for _, record := range records {
		if len(record) != 3 {
			return nil, fmt.Errorf("invalid record length: %d", len(record))
		}

		id, _ := strconv.Atoi(record[0])
		target, _ := strconv.Atoi(record[1])
		ratio, _ := strconv.ParseFloat(record[2], 64)

		problem := KnapsackProblem{
			ID:     id,
			Target: target,
			Ratio:  ratio,
		}

		currentGroup = append(currentGroup, problem)
		if len(currentGroup) == 15 {
			problemsList = append(problemsList, currentGroup)
			currentGroup = make([]KnapsackProblem, 0, 15)
		}
	}

	if len(currentGroup) > 0 {
		problemsList = append(problemsList, currentGroup)
	}

	return problemsList, nil
}

func readBruteTimes(filename string) ([][]float64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) > 0 {
		records = records[1:]
	}

	bruteTimes := make([][]float64, 0)
	currentGroup := make([]float64, 0, 15)

	for _, record := range records {
		if len(record) < 7 {
			continue
		}

		timeMs, err := strconv.ParseFloat(record[6], 64)
		if err != nil {
			return nil, err
		}

		currentGroup = append(currentGroup, timeMs)
		if len(currentGroup) == 15 {
			bruteTimes = append(bruteTimes, currentGroup)
			currentGroup = make([]float64, 0, 15)
		}
	}

	if len(currentGroup) > 0 {
		bruteTimes = append(bruteTimes, currentGroup)
	}

	return bruteTimes, nil
}
