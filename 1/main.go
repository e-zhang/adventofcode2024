package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)

	first := []int{}
	second := []int{}

	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.SplitN(line, " ", 2)
		var f, s int
		_, err := fmt.Sscanf(strings.TrimSpace(tokens[0]), "%d", &f)
		if err != nil {
			panic(err)
		}
		first = append(first, f)
		_, err = fmt.Sscanf(strings.TrimSpace(tokens[1]), "%d", &s)
		if err != nil {
			panic(err)
		}
		second = append(second, s)
	}

	part1(first, second)
	part2(first, second)
}

func part2(first, second []int) {
	s := map[int]int{}

	for _, n := range second {
		if _, ok := s[n]; ok {
			s[n] += 1
		} else {
			s[n] = 1
		}
	}

	total := 0
	for _, n := range first {
		v, ok := s[n]
		if !ok {
			continue
		}
		total += n * v
	}

	fmt.Println(total)
}

func part1(first, second []int) {
	sort.Ints(first)
	sort.Ints(second)

	total := 0
	for i, n := range first {
		d := abs(n - second[i])
		total += d
	}

	fmt.Println(total)
}
