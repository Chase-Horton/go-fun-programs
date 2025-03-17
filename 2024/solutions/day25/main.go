package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/chase-horton/advent-of-code-go/utils"
)

type Message struct {
	isLock bool
	pins   []int
}

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}
func getLockPins(lines []string) []int {
	pins := make([]int, 0, len(lines[0]))
	for x := 0; x < len(lines[0]); x++ {
		currPinHeight := 0
		for y := 1; y < len(lines); y++ {
			if lines[y][x] == '#' {
				currPinHeight++
			} else {
				break
			}
		}
		pins = append(pins, currPinHeight)
	}
	return pins
}
func getKeyPins(lines []string) []int {
	pins := make([]int, 0, len(lines[0]))
	for x := 0; x < len(lines[0]); x++ {
		currPinHeight := 0
		for y := len(lines) - 2; y >= 0; y-- {
			if lines[y][x] == '#' {
				currPinHeight++
			} else {
				break
			}
		}
		pins = append(pins, currPinHeight)
	}
	return pins
}
func parseLockOrKey(lockOrKey string) Message {
	lines := strings.Split(strings.TrimSpace(lockOrKey), "\n")
	if lines[0][0] == '#' {
		return Message{isLock: true, pins: getLockPins(lines)}
	} else {

		return Message{isLock: false, pins: getKeyPins(lines)}
	}
}
func compareLockAndKey(lock []int, key []int) bool {
	if len(lock) != len(key) {
		return false
	}
	for i := 0; i < len(lock); i++ {
		if lock[i]+key[i] > 5 {
			return false
		}
	}
	return true
}
func compareKeyToLocks(key []int, locks [][]int) int {
	matches := 0
	for _, lock := range locks {
		if compareLockAndKey(lock, key) {
			matches++
		}
	}
	return matches
}
func main() {
	defer timer("main")()
	data := strings.Split(utils.ReadFile("day25.txt"), "\n\n")

	locks := [][]int{}
	keys := [][]int{}

	for _, lockOrKey := range data {
		msg := parseLockOrKey(lockOrKey)
		if msg.isLock {
			locks = append(locks, msg.pins)
		} else {
			keys = append(keys, msg.pins)
		}
	}

	sum := 0
	for _, key := range keys {
		sum += compareKeyToLocks(key, locks)
	}

	fmt.Println("1. Total # of unique matches:", sum)
}
