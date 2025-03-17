package main

import (
	"strconv"
	"strings"

	"github.com/chase-horton/advent-of-code-go/utils"
)

func sortList(lst []int) []int {
	for i := 0; i < len(lst)-1; i++ {
		for j := 0; j < len(lst)-i-1; j++ {
			if lst[j] > lst[j+1] {
				lst[j], lst[j+1] = lst[j+1], lst[j]
			}
		}
	}
	return lst
}
func parseInt(lst []int, s string) []int {
	num, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	lst = append(lst, num)
	return lst
}
func main() {
	data := utils.ReadLines("day1.txt")
	leftList, rightList := []int{}, []int{}
	occurencesRight := make(map[int]int)

	for _, line := range data {
		line = strings.TrimSpace(line)
		splitLine := strings.Split(line, "   ")
		left, right := splitLine[0], splitLine[1]
		leftList = parseInt(leftList, left)
		rightList = parseInt(rightList, right)
		occurencesRight[rightList[len(rightList)-1]]++
	}

	leftList, rightList = sortList(leftList), sortList(rightList)

	totalDifference := 0
	totalSimilarity := 0

	for i := 0; i < len(leftList); i++ {
		totalDifference += utils.Abs(leftList[i] - rightList[i])
		totalSimilarity += occurencesRight[leftList[i]] * leftList[i]
	}
	println("1. Total Difference: ", totalDifference)
	println("2. Total Similarity: ", totalSimilarity)
}
