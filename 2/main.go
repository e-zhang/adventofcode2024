package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type report []int

func (r report) IsSafe() bool {
	return gradual(r, 1) || gradual(r, -1)
}

func (r report) IsSafeWithDampener() bool {
	if r.IsSafe() {
		return true
	}

	for i := 0; i < len(r); i++ {
		var n report
		n = append(n, r[:i]...)
		n = append(n, r[i+1:]...)

		if n.IsSafe() {
			return true
		}
	}

	return false
}

func gradual(r report, dir int) bool {
	prev := r[0] * dir
	for i := 1; i < len(r); i++ {
		curr := r[i] * dir
		if curr <= prev {
			return false
		}
		if diff := abs(curr - prev); diff > 3 {
			return false
		}
		prev = curr
	}

	return true
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)

	part1 := 0
	part2 := 0
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, " ")

		var r report
		for _, t := range tokens {
			var n int
			_, err := fmt.Sscanf(t, "%d", &n)
			if err != nil {
				panic(err)
			}
			r = append(r, n)
		}

		if r.IsSafe() {
			part1 += 1
		}

		if r.IsSafeWithDampener() {
			part2 += 1
		}
	}

	fmt.Println("part1", part1)
	fmt.Println("part2", part2)
}
