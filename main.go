package main

import (
	"bufio"
	"fmt"
	"log"
	"maps"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Measurements struct {
	min   float64
	max   float64
	sum   float64
	count int
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: dash [measurements.txt]")
	}

	input := os.Args[1]
	file, err := os.Open(input)
	if err != nil {
		log.Panic("Unable to open ", input)
	}
	defer file.Close()

	m := make(map[string]Measurements)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		pair := strings.Split(line, ";")
		station := pair[0]
		temp, _ := strconv.ParseFloat(pair[1], 64)
		measurements, found := m[station]
		if !found {
			measurements = Measurements{temp, temp, temp, 1}
		} else {
			measurements.min = min(temp, measurements.min)
			measurements.max = max(temp, measurements.max)
			measurements.sum += temp
			measurements.count += 1
		}
		m[station] = measurements
	}
	fmt.Printf("{")
	for i, station := range slices.Sorted(maps.Keys(m)) {
		measurements := m[station]

		if i > 0 {
			fmt.Printf(", ")
		}

		fmt.Printf(
			"%s=%.1f/%.1f/%.1f",
			station,
			measurements.min,
			measurements.sum/float64(measurements.count),
			measurements.max)

	}
	fmt.Printf("}")
}
