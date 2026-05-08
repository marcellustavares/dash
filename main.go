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

		begin := 0
		semi := 0
		for i := 0; i < len(data); i++ {
			switch data[i] {
			case ';':
				semi = i
			case '\n':
				processRow(data[begin:semi], data[semi+1:i], m, &totalStations)
				begin = i + 1
			}
		}

		// incomplete line
		if begin < len(data) {
			leftover = append(leftover[:0], data[begin:]...)
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

func parseFloat(bytes []byte) float64 {
	i := 0
	sign := 1.
	if bytes[i] == '-' {
		sign = -1
		i++
	}
	temp := 0.0
	for ; bytes[i] != '.'; i++ {
		temp = temp*10 + float64(bytes[i]-'0')
	}
	i++
	decimal := float64(bytes[i] - '0')
	temp += decimal / 10.0
	return sign * temp
}

func processRow(stationBytes []byte, temperatureBytes []byte, m []Entry, totalStations *int) {
	hashCode := HashBytes(stationBytes)
	index := int(hashCode % BUCKET_SIZE)
	temp := parseFloat(temperatureBytes)
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
