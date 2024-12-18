package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var VERBOSE bool
var FILE = "input"

var (
	SIZE  = 70
	PART1 = 1024
)

func init() {
	flag.BoolVar(&VERBOSE, "v", false, "print out extra debug info")

	flag.Parse()
	if flag.NArg() > 0 {
		FILE = flag.Arg(0)
		if FILE != "input" {
			SIZE = 6
			PART1 = 12
		}
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
	bytes := []Coord{}
	for scanner.Scan() {
		line := scanner.Text()
		var x, y int
		if _, err := fmt.Sscanf(line, "%d,%d", &x, &y); err != nil {
			panic(err)
		}
		bytes = append(bytes, Coord{x, y})
	}

	start := Coord{0, 0}
	end := Coord{SIZE, SIZE}
	fmt.Println(BFS(start, end, bytes[:PART1]))

	for i := 0; i < len(bytes); i++ {
		s := BFS(start, end, bytes[:i])
		if s < 0 {
			fmt.Println(i, bytes[i-1])
			break
		}
	}
}

type Coord struct {
	x int
	y int
}

func (c Coord) Add(o Coord) Coord {
	return Coord{c.x + o.x, c.y + o.y}
}

type node struct {
	pos   Coord
	steps int
}

func BFS(start, end Coord, bytes []Coord) int {
	q := []node{{start, 0}}

	lut := map[Coord]struct{}{}
	for _, b := range bytes {
		lut[b] = struct{}{}
	}

	seen := map[Coord]struct{}{}

	for len(q) > 0 {
		p := q[0]
		q = q[1:]

		if p.pos == end {
			return p.steps
		}

		if _, ok := seen[p.pos]; ok {
			continue
		}
		seen[p.pos] = struct{}{}

		for _, n := range []Coord{
			{1, 0},
			{-1, 0},
			{0, 1},
			{0, -1},
		} {
			next := p.pos.Add(n)
			if next.x < 0 || next.x > end.x ||
				next.y < 0 || next.y > end.y {
				continue
			}
			if _, ok := lut[next]; ok {
				continue
			}

			q = append(q, node{next, p.steps + 1})
		}
	}

	return -1
}
