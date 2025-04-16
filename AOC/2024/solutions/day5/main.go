package main

import (
	"fmt"
	"strings"

	"github.com/chase-horton/advent-of-code-go/utils"
)

type Rule struct {
	From     int
	To       int
	FromSeen bool
	Finished bool
}

func NewRule(line string) Rule {
	chunks := strings.Split(line, "|")
	return Rule{
		From:     utils.StrToIntorPanic(strings.TrimSpace(chunks[0])),
		To:       utils.StrToIntorPanic(strings.TrimSpace(chunks[1])),
		FromSeen: false,
		Finished: false,
	}
}
func getMiddleValue(s []string) int {
	middle := int(len(s)/2) + 1
	return utils.StrToIntorPanic(s[middle])
}
func resetRules(rules []Rule) {
	for i := range rules {
		rules[i].FromSeen = false
		rules[i].Finished = false
	}
}
func main() {
	data := utils.ReadLines("day5.txt")
	rules := make([]Rule, 0)
	for _, line := range data {
		if line != "\r" {
			break
		}
		rules = append(rules, NewRule(line))
	}
	numSuccess := 0
	for _, line := range data[1177:] {
		resetRules(rules)
		failed := false
		for _, number := range strings.Split(line, ",") {
			num := utils.StrToIntorPanic(strings.TrimSpace(number))
			for i := range rules {
				if rules[i].From == num && !rules[i].FromSeen {
					rules[i].FromSeen = true
				}
				if rules[i].To == num && !rules[i].FromSeen {
					failed = true
					break
				}
				if rules[i].FromSeen && rules[i].To == num {
					rules[i].Finished = true
				}
			}
			if failed {
				break
			}
		}
		for i := range rules {
			if rules[i].FromSeen && !rules[i].Finished {
				failed = true
				break
			}
		}
		if !failed {
			numSuccess += getMiddleValue(strings.Split(line, ","))
		}
	}
	fmt.Printf("Number of successful lines: %d", numSuccess)
}
