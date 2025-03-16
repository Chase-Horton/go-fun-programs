package utils

import (
	"os"
	"strings"
)

func ReadFile(filename string) string {
	data, err := os.ReadFile("./2024/data/" + filename)
	if err != nil {
		panic(err)
	}
	return string(data)
}
func ReadLines(filename string) []string {
	str := ReadFile(filename)
	return strings.Split(str, "\n")
}
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
