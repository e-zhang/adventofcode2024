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
	scanner.Scan()
	line := scanner.Text()
	stones := strings.Split(line, " ")

	fmt.Println(stones)
	sum := 0
	seen := map[step]int{}
	for _, s := range stones {
		sum += BlinkTimes(s, 25, seen)
	}
	fmt.Println(sum, BlinkCount(stones, 25))
	sum = 0
	for _, s := range stones {
		sum += BlinkTimes(s, 75, seen)
	}
	fmt.Println(sum, BlinkCount(stones, 75))
}

type step struct {
	stone string
	i     int
}

func BlinkCount(s []string, iter int) int {
	stones := map[string]int{}
	for _, x := range s {
		stones[x] = 1
	}

	for i := 0; i < iter; i++ {
		next := map[string]int{}
		for x, qty := range stones {
			n := BlinkOnce(x)
			for _, y := range n {
				if v, ok := next[y]; ok {
					next[y] = v + qty
				} else {
					next[y] = qty
				}
			}
		}
		stones = next
	}

	sum := 0
	for _, qty := range stones {
		sum += qty
	}
	return sum
}

func BlinkTimes(stone string, iter int, seen map[step]int) int {
	if v, ok := seen[step{stone, iter}]; ok {
		return v
	}

	next := BlinkOnce(stone)

	if iter == 1 {
		seen[step{stone, iter}] = len(next)
		return len(next)
	}

	sum := 0
	for _, s := range next {
		sum += BlinkTimes(s, iter-1, seen)
	}

	seen[step{stone, iter}] = sum
	return sum
}

func BlinkOnce(s string) []string {
	next := []string{}
	switch {
	case s == "0":
		next = append(next, "1")
	case len(s)%2 == 0:
		next = append(next, s[:len(s)/2])
		x := strings.TrimLeft(s[len(s)/2:], "0")
		if len(x) == 0 {
			x = "0"
		}
		next = append(next, x)
	default:
		var i int
		if _, err := fmt.Sscanf(s, "%d", &i); err != nil {
			panic(err)
		}
		next = append(next, fmt.Sprintf("%d", i*2024))
	}
	return next
}
