package main

import (
	"bufio"
	"container/heap"
	"flag"
	"fmt"
	"os"
)

var (
	VERBOSE bool
	FILE    = "input"
)

var (
	E = Coord{0, 1}
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
	// seen := map[key]int{}
	// chain := path{-1, nil}
	// score, _ := MazeDFS(start, E, end, []Coord{}, grid, 0, seen, &chain)
	// score, tiles := MazeBFS(start, E, end, grid)
	score, tiles := MazeDijkstras(start, E, end, grid)
	fmt.Println(score, tiles)
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

func MazeDFS(pos, dir, end Coord, prevs []Coord, grid []string, score int, seen map[key]int, chain *path) (int, []Coord) {
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

	options := []cache{}

	cw, cwPath := MazeDFS(pos, dir.Rotate(true), end, prevs, grid, score+1000, seen, chain)
	options = append(options, cache{cw, cwPath})
	ccw, ccwPath := MazeDFS(pos, dir.Rotate(false), end, prevs, grid, score+1000, seen, chain)
	options = append(options, cache{ccw, ccwPath})
	next := pos.Add(dir)
	if grid[next.row][next.col] != '#' {
		moveScore, movePath := MazeDFS(next, dir, end, append(prevs, pos), grid, score+1, seen, chain)
		options = append(options, cache{moveScore, movePath})
	}

	best := options[0]
	for _, b := range options[1:] {
		if better(best.score, b.score) {
			best = b
		}
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

type node struct {
	pos   Coord
	dir   Coord
	score int
	path  []Coord
}

type priorityq []node

func (p priorityq) Len() int           { return len(p) }
func (p priorityq) Less(i, j int) bool { return p[i].score < p[j].score }
func (p priorityq) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p *priorityq) Push(x any)        { *p = append(*p, x.(node)) }
func (p *priorityq) Pop() any {
	old := *p
	n := len(old)
	item := old[n-1]
	*p = old[:n-1]
	return item
}

func MazeBFS(pos, dir, end Coord, grid []string) (int, int) {
	q := []node{{pos, dir, 0, nil}}
	seen := map[key]int{}

	best := 0
	tiles := map[Coord]struct{}{}

	for len(q) > 0 {
		curr := q[0]
		q = q[1:]

		if curr.pos == end {
			if best > 0 && curr.score > best {
				continue
			}
			if best == 0 || curr.score < best {
				best = curr.score
				tiles = map[Coord]struct{}{end: struct{}{}}
			}

			for _, c := range curr.path {
				tiles[c] = struct{}{}
			}
			continue
		}

		k := key{curr.pos, curr.dir}
		if v, ok := seen[k]; ok && curr.score > v {
			continue
		}
		seen[k] = curr.score
		path := make([]Coord, len(curr.path))
		copy(path, curr.path)
		for _, n := range []node{
			{curr.pos.Add(curr.dir), curr.dir, curr.score + 1, append(path, curr.pos)},
			{curr.pos, curr.dir.Rotate(true), curr.score + 1000, path},
			{curr.pos, curr.dir.Rotate(false), curr.score + 1000, path},
		} {
			if grid[n.pos.row][n.pos.col] == '#' {
				continue
			}
			q = append(q, n)
		}
	}

	debug(best, len(tiles))
	print(grid, tiles)
	return best, len(tiles)
}

func MazeDijkstras(start, dir, end Coord, grid []string) (int, int) {
	q := priorityq{{start, dir, 0, nil}}
	heap.Init(&q)
	dist := map[key]int{key{start, dir}: 0}
	prevs := map[key][]key{}

	best := 0

	for q.Len() > 0 {
		curr := heap.Pop(&q).(node)

		if curr.pos == end {
			if best == 0 || curr.score < best {
				best = curr.score
			}
			continue
		}

		for _, n := range []node{
			{curr.pos.Add(curr.dir), curr.dir, curr.score + 1, nil},
			{curr.pos, curr.dir.Rotate(true), curr.score + 1000, nil},
			{curr.pos, curr.dir.Rotate(false), curr.score + 1000, nil},
		} {
			if grid[n.pos.row][n.pos.col] == '#' {
				continue
			}

			k := key{n.pos, n.dir}
			v := dist[k]
			if v == 0 || n.score <= v {
				dist[k] = n.score
				if n.score == v {
					prevs[k] = append(prevs[k], key{curr.pos, curr.dir})
				} else {
					prevs[k] = []key{{curr.pos, curr.dir}}
				}
				heap.Push(&q, n)
			}
		}
	}

	tiles := 0
	for k, v := range dist {
		if k.pos == end && v == best {
			tiles += backtrack(start, k.dir, end, grid, prevs)
		}
	}
	return best, tiles
}

func backtrack(start, dir, end Coord, grid []string, prevs map[key][]key) int {
	tiles := map[Coord]struct{}{}
	seen := map[key]struct{}{}

	q := []key{{end, dir}}
	for len(q) > 0 {
		pos := q[0]
		q = q[1:]

		if _, ok := seen[pos]; ok {
			continue
		}

		seen[pos] = struct{}{}
		tiles[pos.pos] = struct{}{}
		if pos.pos == start {
			continue
		}

		q = append(q, prevs[pos]...)
	}

	debug(len(tiles))
	print(grid, tiles)
	return len(tiles)
}
