package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var VERBOSE bool
var FILE = "input"

var ()

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
	patterns := []string{}
	towels := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		if len(patterns) == 0 {
			patterns = strings.Split(line, ", ")
		} else {
			towels = append(towels, line)
		}
	}

	count := 0
	combos := 0
	seen := map[string]int{}
	for _, t := range towels {
		if x := design(t, patterns, seen); x > 0 {
			debug(t)
			combos += x
			count++
		}
	}
	fmt.Println(count, combos)
}

func design(towel string, patterns []string, seen map[string]int) int {
	if len(towel) == 0 {
		return 1
	}

	if v, ok := seen[towel]; ok {
		return v
	}

	combos := 0
	for _, p := range patterns {
		if strings.HasPrefix(towel, p) {
			combos += design(towel[len(p):], patterns, seen)
		}
	}

	seen[towel] = combos
	return combos
}
