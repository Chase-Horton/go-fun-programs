package main

import (
	"regexp"
	"strings"

	"github.com/chase-horton/advent-of-code-go/utils"
)

func calculateSumOfMatches(matches [][]string) int {
	sum := 0
	for _, match := range matches {
		sum += utils.StrToIntorPanic(match[1]) * utils.StrToIntorPanic(match[2])
	}
	return sum
}

func main() {
	data := utils.ReadFile("day3.txt")
	reMull := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)
	reDoDont := regexp.MustCompile(`(?s)don't\(\).*?do\(\)`)

	matches := reMull.FindAllStringSubmatch(data, -1)
	println("1. Sum of multiplications:", calculateSumOfMatches(matches))

	toRemove := reDoDont.FindAllString(data, -1)
	for _, remove := range toRemove {
		data = strings.Replace(data, remove, "", 1)
	}
	matches = reMull.FindAllStringSubmatch(data, -1)
	println("2. Sum of multiplications without don'ts:", calculateSumOfMatches(matches))

}
