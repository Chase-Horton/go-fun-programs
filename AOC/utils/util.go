package utils

import (
	"os"
	"strconv"
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
	str = strings.TrimSpace(str)
	return strings.Split(str, "\n")
}
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
func StrToIntorPanic(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
