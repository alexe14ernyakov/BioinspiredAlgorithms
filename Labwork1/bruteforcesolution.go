package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Item struct {
	Weight int
	Index  int
}

type KnapsackProblem struct {
	ID     int
	Target int
	Ratio  float64
}

type Solution struct {
	AchievedWeight    int
	Combinations      [][]int
	FirstSolutionTime float64
	AllSolutionsTime  float64
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

	var itemsList [][]Item
	for _, record := range records {
		var items []Item
		for i, value := range record {
			weight, err := strconv.Atoi(value)
			if err != nil {
				return nil, err
			}
			items = append(items, Item{Weight: weight, Index: i})
		}
		itemsList = append(itemsList, items)
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

	var problemsList [][]KnapsackProblem
	var currentGroup []KnapsackProblem

	for i, record := range records {
		if len(record) != 3 {
			return nil, fmt.Errorf("invalid number of fields in row %d", i+1)
		}

		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, err
		}

		target, err := strconv.Atoi(record[1])
		if err != nil {
			return nil, err
		}

		ratio, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, err
		}

		problem := KnapsackProblem{
			ID:     id,
			Target: target,
			Ratio:  ratio,
		}

		currentGroup = append(currentGroup, problem)
		if len(currentGroup) == 15 {
			problemsList = append(problemsList, currentGroup)
			currentGroup = nil
		}
	}

	if len(currentGroup) > 0 {
		problemsList = append(problemsList, currentGroup)
	}

	return problemsList, nil
}

func solveKnapsack(items []Item, target int) Solution {
	n := len(items)
	maxWeight := 0
	solutions := make([][]int, 0)

	startAll := time.Now()
	var firstSolutionTime time.Time
	found := false

	for mask := 0; mask < (1 << uint(n)); mask++ {
		currentWeight := 0
		var currentCombination []int

		for i := 0; i < n; i++ {
			if mask&(1<<uint(i)) != 0 {
				currentWeight += items[i].Weight
				currentCombination = append(currentCombination, items[i].Index)
			}

			if currentWeight > target {
				break
			}
		}

		if currentWeight > target {
			continue
		}

		if currentWeight > maxWeight {
			maxWeight = currentWeight
			solutions = [][]int{currentCombination}
			if !found {
				firstSolutionTime = time.Now()
				found = true
				fmt.Printf("Первое решение: вес=%d, комбинация=%v\n",
					maxWeight, currentCombination)
			}
		} else if currentWeight == maxWeight {
			solutions = append(solutions, currentCombination)
		}
	}

	allTime := time.Since(startAll)
	var firstTime time.Duration
	if found {
		firstTime = firstSolutionTime.Sub(startAll)
	}

	return Solution{
		AchievedWeight:    maxWeight,
		Combinations:      solutions,
		FirstSolutionTime: float64(firstTime.Microseconds()) / 1000,
		AllSolutionsTime:  float64(allTime.Microseconds()) / 1000,
	}
}

func main() {
	itemsList, err := readItems("knapsack_vectors.csv")
	if err != nil {
		log.Fatalf("Error reading items: %v", err)
	}

	problemsList, err := readProblems("problems.csv")
	if err != nil {
		log.Fatalf("Error reading problems: %v", err)
	}

	if len(itemsList) != len(problemsList) {
		log.Fatalf("Mismatched data: %d item vectors vs %d problem sets",
			len(itemsList), len(problemsList))
	}

	outputFile, err := os.Create("bruteforce_solutions.csv")
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	writer.Write([]string{
		"VectorID", "ProblemID", "TargetWeight", "AchievedWeight",
		"SolutionsCount", "FirstSolutionTime(ms)", "AllSolutionsTime(ms)",
		"ItemsInSolution", "SolutionIndices",
	})

	for vectorID := 0; vectorID < len(itemsList); vectorID++ {
		items := itemsList[vectorID]
		problems := problemsList[vectorID]

		for _, problem := range problems {
			fmt.Printf("\n--- Вектор %d, Задача %d ---\n",
				vectorID+1, problem.ID)

			solution := solveKnapsack(items, problem.Target)

			solutionsStr := ""
			if len(solution.Combinations) > 0 {
				var solutions []string
				for _, comb := range solution.Combinations {
					solutions = append(solutions, fmt.Sprintf("%v", comb))
				}
				solutionsStr = strings.Join(solutions, "; ")
			}

			err := writer.Write([]string{
				strconv.Itoa(vectorID + 1),
				strconv.Itoa(problem.ID),
				strconv.Itoa(problem.Target),
				strconv.Itoa(solution.AchievedWeight),
				strconv.Itoa(len(solution.Combinations)),
				fmt.Sprintf("%.3f", solution.FirstSolutionTime),
				fmt.Sprintf("%.3f", solution.AllSolutionsTime),
				strconv.Itoa(len(items)),
				solutionsStr,
			})

			if err != nil {
				log.Printf("Error writing result: %v", err)
			}

			fmt.Printf(
				"Достигнуто: %d (Цель:%d), Решений=%d\n"+
					"Время первого решения: %.3f мс\n"+
					"Общее время выполнения: %.3f мс\n",
				solution.AchievedWeight, problem.Target, len(solution.Combinations),
				solution.FirstSolutionTime,
				solution.AllSolutionsTime,
			)
		}
	}

	fmt.Println("\nВсе задачи решены. Результаты сохранены в knapsack_solutions_ms.csv")
}
