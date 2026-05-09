package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"sort"
	"time"
)

const BUCKET_SIZE = 250000

type Entry struct {
	stationBytes []byte
	measurements *Measurements
}

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

	m := make([]Entry, BUCKET_SIZE)
	totalStations := 0

	bufferSize := 4 * 1024 * 1024
	reader := bufio.NewReaderSize(file, bufferSize)

	buffer := make([]byte, bufferSize)
	var leftover []byte

	for {
		n, _ := reader.Read(buffer)

		if n == 0 {
			break
		}

		data := buffer[:n]

		if len(leftover) > 0 {
			tmp := make([]byte, len(leftover)+len(data))
			copy(tmp, leftover)
			copy(tmp[len(leftover):], data)
			data = tmp
			leftover = leftover[:0]
		}

		beginLine := 0

		semi := 0
		hashCode := uint64(offset64)
		sign := 1
		temp := 0
		parsingTemp := false
		for i := 0; i < len(data); i++ {
			b := data[i]

			if !parsingTemp {
				if b == ';' {
					semi = i
					parsingTemp = true
					continue
				}

				hashCode ^= uint64(b)
				hashCode *= prime64
				continue
			}

			switch b {
			case '-':
				sign = -1
			case '.':
				// noop
			case '\n':
				tempFloat := float64(sign*temp) / 10.0
				processRow(data[beginLine:semi], tempFloat, hashCode, m, &totalStations)

				hashCode = offset64
				temp = 0
				sign = 1
				parsingTemp = false

				beginLine = i + 1
			default:
				temp = temp*10 + int(b-'0')
			}
		}

		// incomplete line
		if beginLine < len(data) {
			leftover = append(leftover[:0], data[beginLine:]...)
		}
	}

	sortedEntries := make([]Entry, 0, totalStations)
	for _, entry := range m {
		if entry.stationBytes == nil {
			continue
		}
		sortedEntries = append(sortedEntries, entry)
	}
	sort.Slice(sortedEntries, func(i, j int) bool {
		return string(sortedEntries[i].stationBytes) < string(sortedEntries[j].stationBytes)
	})

	fmt.Printf("{")
	for i, entry := range sortedEntries {
		measurements := entry.measurements

		if i > 0 {
			fmt.Printf(", ")
		}

		fmt.Printf(
			"%s=%.1f/%.1f/%.1f",
			entry.stationBytes,
			measurements.min,
			measurements.sum/float64(measurements.count),
			measurements.max)

	}
	fmt.Printf("}\n")

	elapsed := time.Since(start)
	fmt.Printf("Took %v \n", elapsed)
}

func processRow(stationBytes []byte, temp float64, hashCode uint64, m []Entry, totalStations *int) {
	index := int(hashCode % BUCKET_SIZE)
	for {
		entry := m[index]

		if entry.stationBytes == nil {
			stationCopy := make([]byte, len(stationBytes))
			copy(stationCopy, stationBytes)

			m[index] = Entry{
				stationBytes: stationCopy,
				measurements: &Measurements{temp, temp, temp, 1},
			}
			//fmt.Println(string(stationBytes))
			(*totalStations)++
			break
		}

		eq := bytes.Equal(entry.stationBytes, stationBytes)

		if eq {
			measurements := entry.measurements
			measurements.min = min(temp, measurements.min)
			measurements.max = max(temp, measurements.max)
			measurements.sum += temp
			measurements.count++
			break
		}

		index++
		if index >= BUCKET_SIZE {
			index = 0
		}
	}
}
