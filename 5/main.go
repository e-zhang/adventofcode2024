package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)

	ordering := map[string][]string{}
	pages := [][]string{}

	section := 0
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			section++
			continue
		}

		if section == 0 {
			tokens := strings.Split(line, "|")
			if pages, ok := ordering[tokens[0]]; ok {
				ordering[tokens[0]] = append(pages, tokens[1])
			} else {
				ordering[tokens[0]] = []string{tokens[1]}
			}
		}

		if section == 1 {
			pages = append(pages, strings.Split(line, ","))
		}
	}

	var sum1, sum2 int
	for _, pg := range pages {
		if isInOrder(pg, ordering) {
			var mid int
			if _, err := fmt.Sscanf(pg[len(pg)/2], "%d", &mid); err != nil {
				panic(err)
			}
			sum1 += mid
		} else {
			pg = doSort(pg, ordering)
			fmt.Println(pg, pg[len(pg)/2])
			var mid int
			if _, err := fmt.Sscanf(pg[len(pg)/2], "%d", &mid); err != nil {
				panic(err)
			}
			sum2 += mid
		}
	}

	fmt.Println(sum1, sum2)
}

func doSort(pg []string, ordering map[string][]string) []string {
	sort.Slice(pg, func(i, j int) bool {
		p := pg[i]
		q := pg[j]

		for _, x := range ordering[p] {
			if q == x {
				return true
			}
		}

		return false
	})

	if !isInOrder(pg, ordering) {
		panic(pg)
	}

	return pg
}

func isInOrder(pg []string, ordering map[string][]string) bool {
	for i, p := range pg {
		// everything in the list after is after
		for _, q := range pg[i+1:] {
			isAfter := false
			for _, x := range ordering[p] {
				if q == x {
					isAfter = true
					break
				}
			}

			if !isAfter {
				return false
			}
		}

		// everything in the list before is before
		for _, q := range pg[:i] {
			isAfter := false
			for _, x := range ordering[q] {
				if p == x {
					isAfter = true
					break
				}
			}

			if !isAfter {
				return false
			}
		}
	}

	return true
}
