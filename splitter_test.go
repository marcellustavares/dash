package main

import (
	"os"
	"testing"
)

func TestSplit(t *testing.T) {
	testFiles := []struct {
		name      string
		numChunks int
	}{
		{"data/sample_data.txt", 2},
		// {"data/measurements_100K.txt", 4},
		// {"data/measurements_1M.txt", 8},
		// {"data/measurements_10M.txt", 16},
	}

	for _, tc := range testFiles {
		t.Run(tc.name, func(t *testing.T) {
			testSplitFile(t, tc.name, tc.numChunks)
		})
	}
}

func testSplitFile(t *testing.T, fileName string, numChunks int) {
	t.Helper()

	stat, err := os.Stat(fileName)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}

	fileSize := stat.Size()

	chunks := Split(fileName, numChunks)

	if len(chunks) != numChunks {
		t.Fatalf(
			"expected %d chunks got %d",
			numChunks,
			len(chunks),
		)
	}

	f, err := os.Open(fileName)
	if err != nil {
		t.Fatalf("open failed: %v", err)
	}
	defer f.Close()

	var totalBytes int64

	for i, chunk := range chunks {
		if chunk.offset > chunk.end {
			t.Fatalf(
				"invalid chunk %d: offset=%d end=%d",
				i,
				chunk.offset,
				chunk.end,
			)
		}

		chunkSize := chunk.end - chunk.offset
		totalBytes += chunkSize

		// assert chunk ends with '\n'
		buf := make([]byte, 1)
		_, err := f.ReadAt(buf, chunk.end-1)
		if err != nil {
			t.Fatalf(
				"failed reading chunk end byte: %v",
				err,
			)
		}

		if buf[0] != '\n' {
			t.Fatalf(
				"chunk %d does not end with newline: got %q",
				i,
				buf[0],
			)
		}

		// assert that chunk should be contiguous
		if i > 0 {
			prev := chunks[i-1]

			if prev.end != chunk.offset {
				t.Fatalf(
					"gap/overlap between chunks %d and %d: prev.end=%d current.offset=%d",
					i-1,
					i,
					prev.end,
					chunk.offset,
				)
			}
		}
	}

	// assert sum of chunk bytes should equal file size
	if totalBytes != fileSize {
		t.Fatalf(
			"total bytes mismatch: got %d want %d",
			totalBytes,
			fileSize,
		)
	}
}
