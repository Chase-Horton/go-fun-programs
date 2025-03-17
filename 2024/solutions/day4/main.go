package main

import (
	"github.com/chase-horton/advent-of-code-go/utils"
)

type WordSearch struct {
	wordMap [][]rune
}

func NewWordSearch(data []string) *WordSearch {
	wordMap := [][]rune{}
	for _, row := range data {
		wordMap = append(wordMap, []rune(row))
	}
	return &WordSearch{
		wordMap: wordMap,
	}
}

func (w *WordSearch) checkWordMatch(x, y, d1, d2 int) bool {
	if len(w.wordMap[0]) <= x+d1*3 || x+d1*3 < 0 || y+d2*3 >= len(w.wordMap) || y+d2*3 < 0 {
		return false
	}
	if w.wordMap[y][x] == 'X' && w.wordMap[y+d2][x+d1] == 'M' && w.wordMap[y+d2*2][x+d1*2] == 'A' && w.wordMap[y+d2*3][x+d1*3] == 'S' {
		return true
	}
	return false
}
func (w *WordSearch) checkWordMatches(x, y int) int {
	sum := 0
	for d1 := -1; d1 < 2; d1++ {
		for d2 := -1; d2 < 2; d2++ {
			if w.checkWordMatch(x, y, d1, d2) {
				sum++
			}
		}
	}
	return sum
}
func main() {
	data := utils.ReadLines("day4.txt")
	wordSearch := NewWordSearch(data)
	s := 0
	for y := 0; y < len(wordSearch.wordMap); y++ {
		for x := 0; x < len(wordSearch.wordMap[0]); x++ {
			s += wordSearch.checkWordMatches(x, y)
		}
	}
	println("1. Number of matches:", s)
}
