package main

import (
	"bytes"
	"io"
	"log"
	"os"
)

type Chunk struct {
	offset int64
	end    int64
}

func Split(fileName string, numChunks int) []Chunk {
	f, _ := os.Open(fileName)
	defer f.Close()

	stat, _ := f.Stat()

	chunkSize := stat.Size() / int64(numChunks)

	buffer := make([]byte, 106) // max line length
	chunks := make([]Chunk, 0, numChunks)

	offset := int64(0)

	for i := 0; i < numChunks-1; i++ {
		seekOffset := offset + chunkSize

		_, err := f.Seek(seekOffset, io.SeekStart)
		if err != nil {
			log.Panic(err)
		}

		n, _ := io.ReadFull(f, buffer)

		c := buffer[:n]

		nlIndex := bytes.IndexByte(c, '\n')
		if nlIndex < 0 {
			log.Panic("newline not found")
		}

		// +1 because end is exclusive
		chunkEnd := offset + chunkSize + int64(nlIndex) + 1

		chunks = append(chunks, Chunk{
			offset: offset,
			end:    chunkEnd,
		})

		offset = chunkEnd
	}

	chunks = append(chunks, Chunk{
		offset: offset,
		end:    stat.Size(),
	})

	return chunks
}
