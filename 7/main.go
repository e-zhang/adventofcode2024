package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)

	var sum1, sum2 int
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, ":")
		var total int
		if _, err := fmt.Sscanf(tokens[0], "%d", &total); err != nil {
			panic(err)
		}

		tokens = strings.Split(strings.TrimSpace(tokens[1]), " ")
		operands := make([]int, len(tokens))
		for i, t := range tokens {
			var x int
			if _, err := fmt.Sscanf(t, "%d", &x); err != nil {
				panic(err)
			}
			operands[i] = x
		}

		running := operands[0]
		if part1(running, operands[1:], total) {
			fmt.Println("part1", total)
			sum1 += total
		}

		if part2(running, operands[1:], total) {
			fmt.Println("part2", total)
			sum2 += total
		}
	}

	fmt.Println(sum1, sum2)
}

func part1(running int, remaining []int, goal int) bool {
	if len(remaining) == 0 {
		return goal == running
	}

	next := remaining[0]
	return part1(running+next, remaining[1:], goal) ||
		part1(running*next, remaining[1:], goal)
}

func part2(running int, remaining []int, goal int) bool {
	if len(remaining) == 0 {
		return goal == running
	}

	next := remaining[0]
	return part2(running+next, remaining[1:], goal) ||
		part2(running*next, remaining[1:], goal) ||
		part2(concat(running, next), remaining[1:], goal)
}

func concat(x, y int) int {
	tmp := y
	for tmp > 0 {
		x *= 10
		tmp /= 10
	}

	return x + y
}
