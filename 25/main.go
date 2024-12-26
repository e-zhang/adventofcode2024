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
	keys := [][]int{}
	locks := [][]int{}

	var curr []string

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			k, l := parse(curr)
			if k != nil {
				keys = append(keys, k)
			}
			if l != nil {
				locks = append(locks, l)
			}
			curr = nil
			continue
		}
		curr = append(curr, line)
	}

	if curr != nil {
		panic(curr)
	}

	debug(keys)
	debug(locks)

	fits := 0
	for _, k := range keys {
		for _, l := range locks {
			fit := true
			for i := 0; i < len(k); i++ {
				if k[i]+l[i] > 5 {
					fit = false
					break
				}
			}

			if fit {
				fits++
			}
		}
	}

	fmt.Println(fits)
}

func parse(curr []string) ([]int, []int) {
	schematic := make([]int, 5)

	var isKey bool
	if curr[0] == "#####" {
		curr = curr[1:]
		isKey = true
	} else {
		curr = curr[:len(curr)-1]
		isKey = false
	}

	for _, row := range curr {
		for i, col := range row {
			if col == '#' {
				schematic[i] += 1
			}
		}
	}

	if isKey {
		return schematic, nil
	}
	return nil, schematic

}
