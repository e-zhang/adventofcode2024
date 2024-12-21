package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var VERBOSE bool
var FILE = "input"

var numericKeypad = []string{
	"789",
	"456",
	"123",
	" 0A",
}

var directionalKeypad = []string{
	" ^A",
	"<v>",
}

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

type cacheKey struct {
	start Coord
	end   Coord
	robot int
}

func main() {
	f, err := os.Open(FILE)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)
	codes := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		codes = append(codes, line)
	}

	complexity1 := 0
	complexity2 := 0
	for _, code := range codes {
		debug(code)
		shortest1 := 0
		shortest2 := 0
		start := Coord{3, 2}

		cache := map[cacheKey]int{}
		for _, c := range code {
			moves := moveKeypad(string(c), start, numericKeypad)
			shortest1 += doDirectional(moves, 2, cache)
			shortest2 += doDirectional(moves, 25, cache)
			start = getKey(c, numericKeypad)
		}
		debug(code, shortest1, shortest2, toNumber(code))
		complexity1 += shortest1 * toNumber(code)
		complexity2 += shortest2 * toNumber(code)
	}

	fmt.Println(complexity1, complexity2)
}

func minLength(x []string) int {
	l := 0
	for _, i := range x {
		if l == 0 || len(i) < l {
			l = len(i)
		}
	}
	return l
}

func toNumber(code string) int {
	var num int
	for _, c := range code {
		if c >= '0' && c <= '9' {
			num = num*10 + int(c-'0')
		}
	}
	return num
}

type Coord struct {
	row int
	col int
}

func (c Coord) Add(o Coord) Coord {
	return Coord{c.row + o.row, c.col + o.col}
}

func (c Coord) ToDirection() string {
	switch c {
	case Coord{-1, 0}:
		return "^"
	case Coord{1, 0}:
		return "v"
	case Coord{0, -1}:
		return "<"
	case Coord{0, 1}:
		return ">"
	}

	panic(c)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func getKey(key rune, keypad []string) Coord {
	for r, row := range keypad {
		for c, col := range row {
			if col == key {
				return Coord{r, c}
			}
		}
	}
	debug(key, keypad)
	panic(string(key))
}

func moveKeypad(code string, start Coord, keypad []string) []string {
	curr := start
	paths := []string{""}

	for _, x := range code {
		k := getKey(x, keypad)
		next := move(curr, k, keypad)
		// debug(next, curr, Coord{r, c}, string(col))
		newPaths := []string{}
		for _, n := range next {
			for _, p := range paths {
				newPaths = append(newPaths, p+n+"A")
			}
		}
		paths = newPaths
		curr = k
	}

	// debug(paths)
	return paths
}

func move(start, end Coord, keypad []string) []string {
	type node struct {
		pos  Coord
		path string
	}

	q := []node{{start, ""}}

	diff := Coord{end.row - start.row, end.col - start.col}

	neighbors := []Coord{}
	if diff.col != 0 {
		dc := Coord{0, (diff.col) / abs(diff.col)}
		neighbors = append(neighbors, dc)
	}
	if diff.row != 0 {
		dr := Coord{(diff.row) / abs(diff.row), 0}
		neighbors = append(neighbors, dr)
	}

	paths := []string{}

	for len(q) > 0 {
		curr := q[0]
		q = q[1:]

		if curr.pos == end {
			paths = append(paths, curr.path)
			// break
			continue
		}

		for _, n := range neighbors {
			next := curr.pos.Add(n)
			if next.row < 0 || next.row >= len(keypad) || next.col < 0 || next.col >= len(keypad[next.row]) {
				continue
			}

			if keypad[next.row][next.col] == ' ' {
				continue
			}

			q = append(q, node{next, curr.path + n.ToDirection()})
		}
	}

	return paths
}

func doDirectional(moves []string, robots int, cache map[cacheKey]int) int {
	if robots == 0 {
		return minLength(moves)
	}

	shortest := 0
	for _, m := range moves {
		start := Coord{0, 2}
		path := 0
		for _, s := range m {
			nextMoves := moveKeypad(string(s), start, directionalKeypad)
			end := getKey(s, directionalKeypad)
			k := cacheKey{start, end, robots}
			if v, ok := cache[k]; ok {
				path += v
			} else {
				seq := doDirectional(nextMoves, robots-1, cache)
				path += seq
				cache[k] = seq
			}
			start = end
		}

		if shortest == 0 || path < shortest {
			shortest = path
		}
	}

	return shortest
}
