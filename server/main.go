// main.go
package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	// data     AppData
	dataFile = "financial_data.json"
	mu       sync.RWMutex
)

// Filter function using generics
func Filter[T any](slice []T, predicate func(T) bool) []T {
	var result []T
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

func pseudo_uuid() (uuid string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	uuid = fmt.Sprintf("%X", b[0:4])

	return
}

func printTimeSlice(ts []time.Time) {
	for _, t := range ts {
		fmt.Println(t)
	}
}

func main() {
	// Get the directory of the currently running executable
	// This approach is generally robust across different OSes.
	ex, err := os.Executable()
	if err != nil {
		panic(fmt.Sprintf("Failed to get executable path: %v", err))
	}
	execDir := filepath.Dir(ex)

	// Construct the full path to your data file
	dataFile = filepath.Join(execDir, dataFile)
	fmt.Printf("Data file path: %s\n", dataFile) // For debugging

	setupHttpEndpoints()
}
