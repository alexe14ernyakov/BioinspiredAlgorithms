package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	minItems   = 2
	maxItems   = 12
	maxTries   = 1000
	numTasks   = 15
	totalItems = 24
)

func main() {
	vectors, err := readVectorsFromCSV("knapsack_vectors.csv")
	if err != nil {
		log.Println("Ошибка при чтении файла:", err)
		return
	}

	rand.Seed(time.Now().UnixNano())

	for i, vector := range vectors {
		fmt.Printf("Вектор %d: %v\n", i+1, vector)
		for taskNum := 1; taskNum <= numTasks; taskNum++ {
			target, selectedItems := generateTask(vector)
			if target == -1 {
				fmt.Printf("  Задача %d: Не удалось найти подходящий target_weight\n", taskNum)
				continue
			}
			selectedItemsCount := len(selectedItems)
			percentage := float64(selectedItemsCount) / float64(totalItems)
			fmt.Printf("  Задача %d: Целевой вес = %d, Выбранные предметы: %v, Доля: %.2f\n",
				taskNum, target, selectedItems, percentage)
		}
		fmt.Println()
	}
}

func readVectorsFromCSV(path string) ([][]int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var vectors [][]int
	for _, record := range records {
		var vec []int
		for _, val := range record {
			num, err := strconv.Atoi(val)
			if err != nil {
				log.Println("Ошибка при преобразовании значения в целое число:", err)
				continue
			}
			vec = append(vec, num)
		}
		vectors = append(vectors, vec)
	}
	return vectors, nil
}

func generateTask(weights []int) (int, []int) {
	n := len(weights)

	numSelectedItems := rand.Intn(maxItems-minItems+1) + minItems

	selectedItems := rand.Perm(n)[:numSelectedItems]

	var totalWeight int
	for _, idx := range selectedItems {
		totalWeight += weights[idx]
	}

	return totalWeight, selectedItems
}
