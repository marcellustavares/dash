package main

import (
	"bufio"
	"fmt"
	"log"
	"maps"
	"os"
	"runtime/pprof"
	"slices"
	"strconv"
	"strings"
	"time"
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

	// cpu profiling
	prof, err := os.Create("cpu.prof")
	if err != nil {
		log.Panic("Unable to create cpu.prof file ", err)
	}
	defer prof.Close()
	pprof.StartCPUProfile(prof)
	defer pprof.StopCPUProfile()

	input := os.Args[1]
	file, err := os.Open(input)
	if err != nil {
		log.Panic("Unable to open ", input)
	}
	defer file.Close()

	start := time.Now()

	m := make(map[string]Measurements)

	reader := bufio.NewReader(file)
	for {
		line, _ := reader.ReadString('\n')
		if len(line) == 0 {
			break
		}
		pair := strings.Split(line[:len(line)-1], ";")
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
	fmt.Printf("}\n")

	elapsed := time.Since(start)
	fmt.Printf("Took %v \n", elapsed)
}
