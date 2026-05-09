package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"maps"
	"os"
	"runtime"
	"runtime/pprof"
	"slices"
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

	start := time.Now()

	input := os.Args[1]

	workers := runtime.NumCPU()
	chunks := Split(input, workers)
	entriesChannel := make(chan []Entry, workers)

	for _, chunk := range chunks {
		go processChunk(input, chunk, entriesChannel)
	}

	totalStations := 0
	m := make(map[string]*Measurements)
	for i := 0; i < workers; i++ {
		entries := <-entriesChannel
		for _, entry := range entries {
			station := string(entry.stationBytes) // TODO
			measurements, found := m[station]

			if !found {
				m[station] = entry.measurements
				totalStations++
				continue
			}

			measurements.min = min(measurements.min, entry.measurements.min)
			measurements.max = max(measurements.max, entry.measurements.max)
			measurements.sum += entry.measurements.sum
			measurements.count += entry.measurements.count
		}
	}

	// sortedKeys := slices.Sorted(maps.Keys(m))
	// sortedEntries := sortEntries(m)

	printOutput(m)

	elapsed := time.Since(start)
	fmt.Printf("Took %v \n", elapsed)
}

func processChunk(input string, chunk Chunk, resultsChannel chan []Entry) {
	file, err := os.Open(input)
	if err != nil {
		log.Panic("Unable to open ", input)
	}
	defer file.Close()

	m := make([]Entry, BUCKET_SIZE)
	totalStations := 0

	bufferSize := 4 * 1024 * 1024

	offset := chunk.offset
	end := chunk.end

	buffer := make([]byte, bufferSize)
	var leftover []byte

	for offset < end {
		n, _ := file.ReadAt(buffer, offset)

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

		offset += int64(n)
	}

	// emit worker results
	results := make([]Entry, 0, totalStations)
	for _, entry := range m {
		if entry.stationBytes == nil {
			continue
		}
		results = append(results, entry)
	}
	resultsChannel <- results
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

func printOutput(m map[string]*Measurements) {
	writer := bufio.NewWriter(os.Stdout)

	fmt.Fprint(writer, "{")
	for i, station := range slices.Sorted(maps.Keys(m)) {
		measurements := m[station]

		if i > 0 {
			fmt.Fprint(writer, ", ")
		}

		fmt.Fprintf(
			writer,
			"%s=%.1f/%.1f/%.1f",
			station,
			measurements.min,
			measurements.sum/float64(measurements.count),
			measurements.max)

	}
	fmt.Fprint(writer, "}\n")
	writer.Flush()
}
