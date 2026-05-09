package main

const offset64 = 14695981039346656037
const prime64 = 1099511628211

func HashBytes(b []byte) uint64 {
	var hash uint64 = offset64
	for i := 0; i < len(b); i++ {
		hash ^= uint64(b[i])
		hash *= prime64
	}
	return hash
}
