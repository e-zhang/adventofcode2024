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
	fmt.Println(sum)
	sum = 0
	for _, s := range stones {
		sum += BlinkTimes(s, 75, seen)
	}
	fmt.Println(sum)
}

type step struct {
	stone string
	i     int
}

func BlinkTimes(stone string, iter int, seen map[step]int) int {
	if v, ok := seen[step{stone, iter}]; ok {
		return v
	}

	next := BlinkOnce([]string{stone})

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

func BlinkOnce(stones []string) []string {
	next := []string{}
	for _, s := range stones {
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
	}
	return next
}
