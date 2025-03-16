package main

import (
	"strings"

	"github.com/chase-horton/advent-of-code-go/utils"
)

func parseLine(line string) []int {
	nums := strings.Split(strings.TrimSpace(line), " ")
	out := []int{}
	for _, num := range nums {
		out = append(out, utils.StrToIntPanic(num))
	}
	return out
}
func evaluateLine(line []int) int {
	problems := 0
	lastNum := line[0]
	descending := line[0] > line[1]
	for _, num := range line[1:] {
		if num > lastNum+3 || num < lastNum-3 {
			problems++
		} else if num == lastNum {
			problems++
		} else if (num > lastNum && descending) || (num < lastNum && !descending) {
			problems++
		}
		lastNum = num
	}
	return problems
}
func main() {
	data := utils.ReadLines("day2.txt")
	safeReports := 0
	safeReports2 := 0

	for _, row := range data {
		line := parseLine(row)
		if evaluateLine(line) == 0 {
			safeReports++
		}
		if evaluateLine(line) == 1 {
			safeReports2++
		}
	}
	println("1. Safe Reports:", safeReports)
	println("1. Safe Reports 2:", safeReports+safeReports2)

}
