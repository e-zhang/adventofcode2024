package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var (
	VERBOSE bool
	FILE    = "input"
)

var (
	N = Coord{-1, 0}
	E = Coord{0, 1}
	S = Coord{1, 0}
	W = Coord{0, -1}
)

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
	grid := []string{}
	var start, end Coord
	for scanner.Scan() {
		line := scanner.Text()
		for i, c := range line {
			switch c {
			case 'S':
				start = Coord{len(grid), i}
			case 'E':
				end = Coord{len(grid), i}
			}
		}
		grid = append(grid, line)
	}

	// init := Coord{-1, -1}
	seen := map[key]int{}
	remaining := map[key]cache{}
	chain := path{-1, nil}
	score, _ := Maze(start, E, end, []Coord{}, grid, 0, seen, remaining, &chain)
	fmt.Println(score)

	// path := map[Coord]struct{}{end: struct{}{}}
	// backtrack(end, init, chain, path)
	debug(len(chain.tiles) + 1)
	print(grid, chain.tiles)
}

func backtrack(pos, end Coord, chain map[Coord]map[Coord]struct{}, path map[Coord]struct{}) {
	if pos == end {
		return
	}

	prevs := chain[pos]
	debug(pos, prevs)
	for p := range prevs {
		if _, ok := path[p]; ok {
			continue
		}
		path[p] = struct{}{}
		backtrack(p, end, chain, path)
	}

}

func print(grid []string, path map[Coord]struct{}) {
	if !VERBOSE {
		return
	}

	for r, row := range grid {
		for c, cell := range row {
			if _, ok := path[Coord{r, c}]; ok {
				fmt.Printf("O")
			} else {
				fmt.Printf(string(cell))
			}
		}
		fmt.Println()
	}
}

type key struct {
	pos Coord
	dir Coord
}

type path struct {
	score int
	tiles map[Coord]struct{}
}

func better(x, y int) bool {
	if x < 0 && y < 0 {
		return false
	}

	if y < 0 {
		return false
	}

	if x < 0 {
		return true
	}

	if y < x {
		return true
	}
	return false
}

type cache struct {
	score int
	path  []Coord
}

func Maze(pos, dir, end Coord, prevs []Coord, grid []string, score int, seen map[key]int, remaining map[key]cache, chain *path) (int, []Coord) {
	k := key{pos, dir}
	v, ok := seen[k]
	if ok {
		if score > v {
			return -1, nil
		}
	}

	seen[k] = score
	if pos == end {
		// debug(chain.score, score, prevs)
		if chain.score < 0 || score <= chain.score {
			if chain.score < 0 || score < chain.score {
				chain.tiles = map[Coord]struct{}{}
				chain.score = score
			}
			for _, c := range prevs {
				chain.tiles[c] = struct{}{}
			}
		}
		return score, prevs
	}

	if v, ok := remaining[k]; ok {
		return Maze(end, dir, end, append(prevs, v.path...), grid, score+v.score, seen, remaining, chain)
	}

	options := []cache{}

	cw, cwPath := Maze(pos, dir.Rotate(true), end, prevs, grid, score+1000, seen, remaining, chain)
	options = append(options, cache{cw, cwPath})
	ccw, ccwPath := Maze(pos, dir.Rotate(false), end, prevs, grid, score+1000, seen, remaining, chain)
	options = append(options, cache{ccw, ccwPath})
	next := pos.Add(dir)
	if grid[next.row][next.col] != '#' {
		moveScore, movePath := Maze(next, dir, end, append(prevs, pos), grid, score+1, seen, remaining, chain)
		options = append(options, cache{moveScore, movePath})
	}

	best := options[0]
	for _, b := range options[1:] {
		if better(best.score, b.score) {
			best = b
		}
	}

	// debug(pos, best.score, score, options)
	if best.score >= score {
		// remaining[k] = cache{best.score - score, best.path}
	}
	return best.score, best.path
}

type Coord struct {
	row int
	col int
}

func (c Coord) Add(o Coord) Coord {
	return Coord{c.row + o.row, c.col + o.col}
}

func (c Coord) Rotate(clockwise bool) Coord {
	if clockwise {
		return Coord{-c.col, c.row}
	} else {
		return Coord{c.col, -c.row}
	}
}
