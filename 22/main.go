package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var VERBOSE bool
var FILE = "input"

func init() {
	flag.BoolVar(&VERBOSE, "v", false, "print out extra debug info")

	flag.Parse()
	if flag.NArg() > 0 {
		FILE = flag.Arg(0)
	}
}

func debug(a ...any) {
	if VERBOSE {
		fmt.Println(a...)
	}
}

func main() {
	f, err := os.Open(FILE)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)
	seeds := []int{}
	for scanner.Scan() {
		line := scanner.Text()
		var s int
		if _, err := fmt.Sscanf(line, "%d", &s); err != nil {
			panic(err)
		}
		seeds = append(seeds, s)
	}

	// sum := 0
	// for _, s := range seeds {
	// 	seed := s
	// 	for i := 0; i < 2000; i++ {
	// 		s = Evolve(s)
	// 	}
	// 	debug(seed, s)
	// 	sum += s
	// }
	// fmt.Println(sum)

	// prices := [][]int{}
	prices := []map[string]int{}
	for _, s := range seeds {
		buyer := map[string]int{}
		last := s % 10
		changes := []int{}
		for i := 0; i < 2000; i++ {
			s = Evolve(s)
			if len(changes) == 4 {
				changes = changes[1:]
			}
			digit := s % 10
			changes = append(changes, digit-last)

			k := toString(changes)
			if _, ok := buyer[k]; !ok {
				buyer[k] = digit
			}
			last = digit
			if i >= 4 && len(changes) != 4 {
				panic(changes)
			}
		}
		prices = append(prices, buyer)
	}

	max := 0
	seen := map[string]struct{}{}
	for _, b := range prices {
		for k := range b {
			if _, ok := seen[k]; ok {
				continue
			}
			seen[k] = struct{}{}
			bananas := 0
			for i, p := range prices {
				bananas += p[k]

				if (len(p)-i)*9 < max {
					break
				}
			}
			if bananas > max {
				debug(bananas, k)
				max = bananas
			}
		}
	}
	fmt.Println(max)
}

func toString(x []int) string {
	s := ""
	for i, v := range x {
		if i > 0 {
			s += ","
		}
		s += fmt.Sprintf("%d", v)
	}

	return s
}

func calcChanges(prices []int) []int {
	changes := []int{}
	for i := 1; i < 4; i++ {
		changes = append(changes, prices[i]-prices[i-1])
	}
	return changes
}

func cmp(x, y []int) bool {
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

func checkSequences(idx int, prices [][]int) int {
	buyer := prices[idx]
	changes := calcChanges(buyer)
	max := 0
	for i := 4; i < len(buyer); i++ {
		changes = append(changes[len(changes)-3:], buyer[i]-buyer[i-1])
		sum := 0
		for _, p := range prices {
			pchanges := calcChanges(p)
			for j := 4; j < len(p); j++ {
				pchanges = append(pchanges[len(pchanges)-3:], p[j]-p[j-1])
				if cmp(changes, pchanges) {
					sum += p[j]
					break
				}
			}
		}
		// sum += buyer[i]
		// debug("buyer", changes, sum)
		if sum > max {
			max = sum
		}
	}
	return max
}

func Evolve(secret int) int {
	secret = prune(mix(secret*64, secret))
	secret = prune(mix(secret/32, secret))
	secret = prune(mix(secret*2048, secret))
	return secret
}

func mix(v, secret int) int {
	return v ^ secret
}

func prune(v int) int {
	return v % 16777216
}
