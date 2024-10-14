package main

import (
	"fmt"
	"time"
)

func printMultiProgressBar(progress []int, width int) {
	for i := range progress {
		printProgressBar(i, progress[i], width)
	}
	// Move the cursor back up by the number of progress bars printed
	fmt.Printf("\033[%dA", len(progress))
}

func printProgressBar(index int, percent int, width int) {
	// Calculate the number of "=" and spaces for the progress bar
	progress := int(float64(width) * (float64(percent) / 100.0))
	remaining := width - progress

	// Build the progress bar string
	bar := fmt.Sprintf("[%s%s] %d%%",
		string(repeat('=', progress)), string(repeat(' ', remaining)), percent)

	// Print the progress bar on its line
	fmt.Printf("\rTask %d: %s\n", index+1, bar)
}

func repeat(char rune, count int) string {
	// Create a string with repeated characters
	return string(make([]rune, count, count))
}

func main() {
	width := 50
	tasks := 3 // Number of tasks (and progress bars)
	progress := make([]int, tasks)

	for completed := 0; completed <= 100; completed++ {
		for i := range progress {
			// Simulate each task progressing at a different rate
			progress[i] = (completed + i*10) % 101
		}

		// Print all progress bars
		printMultiProgressBar(progress, width)

		// Simulate work with a sleep
		time.Sleep(100 * time.Millisecond)
	}

	// Print a newline after the progress bars complete
	fmt.Println()
}
