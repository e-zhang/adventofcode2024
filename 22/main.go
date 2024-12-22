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

	sequences := map[string]int{}
	part1 := 0
	for _, s := range seeds {
		part1 += doBuyer(s, sequences)
	}

	part2 := 0
	for k, v := range sequences {
		if v > part2 {
			debug(k, v)
			part2 = v
		}
	}

	fmt.Println(part1, part2)
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

func doBuyer(seed int, sequences map[string]int) int {
	buyer := map[string]int{}
	s := seed
	last := s % 10
	changes := []int{}
	for i := 0; i < 2000; i++ {
		s = Evolve(s)
		digit := s % 10

		// update changes tracking
		if len(changes) == 4 {
			changes = changes[1:]
		}
		changes = append(changes, digit-last)
		if i >= 4 && len(changes) != 4 {
			panic(changes)
		}

		k := toString(changes)
		if _, ok := buyer[k]; !ok {
			buyer[k] = digit
			sequences[k] += digit
		}
		last = digit
	}
	debug(seed, s)
	return s
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
