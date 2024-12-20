package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var VERBOSE bool
var FILE = "input"

var CUTOFF = 100

func init() {
	flag.BoolVar(&VERBOSE, "v", false, "print out extra debug info")

	flag.Parse()
	if flag.NArg() > 0 {
		FILE = flag.Arg(0)
		if FILE != "input" {
			CUTOFF = 1
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
	debug(start, end)

	// steps := race(start, end, cheat{}, grid, 0)
	// fmt.Println(steps)
	// saves := map[int]int{}
	// for r, row := range grid {
	// 	for c, _ := range row {
	// 		if grid[r][c] == '#' {
	// 			continue
	// 		}
	// 		cheats := findCheats(Coord{r, c}, grid)
	// 		for _, cheat := range cheats {
	// 			s := race(start, end, cheat, grid, steps-CUTOFF)
	// 			if s > 0 {
	// 				debug(cheat, s, steps-s)
	// 				count := saves[steps-s]
	// 				saves[steps-s] = count + 1
	// 			}
	// 		}

	// 	}
	// }

	d := calculateDistances(start, grid)
	part1 := calculateDistancesWithCheats(grid, d, 2, CUTOFF)
	part2 := calculateDistancesWithCheats(grid, d, 20, CUTOFF)

	sum := 0
	debug(part1)
	for _, v := range part1 {
		sum += v
	}
	fmt.Println(sum)

	sum = 0
	debug(part2)
	for _, v := range part2 {
		sum += v
	}
	fmt.Println(sum)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func calculateDistancesWithCheats(grid []string, distances map[Coord]int, cheat int, cutoff int) map[int]int {
	savings := map[int]int{}

	for d, steps := range distances {
		for r := -cheat; r <= cheat; r++ {
			for c := -cheat; c <= cheat; c++ {
				n := d.Add(Coord{r, c})
				if n.row < 0 || n.col < 0 || n.row >= len(grid) || n.col >= len(grid[n.row]) {
					continue
				}

				if grid[n.row][n.col] == '#' {
					continue
				}

				if steps-distances[n] < cutoff {
					continue
				}

				diff := abs(r) + abs(c)
				if diff > cheat {
					continue
				}

				save := steps - distances[n] - diff
				if save < cutoff {
					continue
				}
				savings[save] += 1
			}
		}
	}

	return savings
}

func calculateDistances(start Coord, grid []string) map[Coord]int {
	type node struct {
		pos   Coord
		steps int
	}
	q := []node{{start, 0}}

	distances := map[Coord]int{}

	for len(q) > 0 {
		p := q[0]
		q = q[1:]

		if _, ok := distances[p.pos]; ok {
			continue
		}
		distances[p.pos] = p.steps

		for _, n := range []Coord{
			{1, 0},
			{-1, 0},
			{0, 1},
			{0, -1},
		} {
			next := p.pos.Add(n)
			if next.row < 0 || next.row >= len(grid) ||
				next.col < 0 || next.col >= len(grid[next.row]) {
				continue
			}

			if grid[next.row][next.col] == '#' {
				continue
			}

			q = append(q, node{next, p.steps + 1})
		}
	}
	return distances
}

type Coord struct {
	row int
	col int
}

func (c Coord) Add(o Coord) Coord {
	return Coord{c.row + o.row, c.col + o.col}
}

/*** OLD CODE that didnt work for part2
type cheat struct {
	start Coord
	end   Coord
	steps int
}

type node struct {
	pos   Coord
	cheat cheat
	steps int
}

// find cheats incorrectly assumed that a cheat had to be strictly wall-to-wall, not just any N spaces through a wall with mix
func findCheats(start Coord, grid []string) []cheat {
	q := []node{{start, cheat{start, start, 0}, 0}}

	m := map[Coord]cheat{}
	seen := map[Coord]struct{}{}

	for len(q) > 0 {
		p := q[0]
		q = q[1:]

		if p.steps > 20 {
			continue
		}

		if p.steps > 0 {
			if grid[p.pos.row][p.pos.col] != '#' {
				if c, ok := m[p.pos]; !ok || p.cheat.steps < c.steps {
					m[p.pos] = p.cheat
				}
				continue
			}
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
			if next.row < 0 || next.row >= len(grid) ||
				next.col < 0 || next.col >= len(grid[next.row]) {
				continue
			}

			q = append(q, node{next, cheat{start, next, p.cheat.steps + 1}, p.steps + 1})
		}
	}

	cheats := []cheat{}
	for _, v := range m {
		cheats = append(cheats, v)
	}
	return cheats
}

func race(start, end Coord, c cheat, grid []string, cutoff int) int {
	q := []node{{start, cheat{}, 0}}

	seen := map[Coord]struct{}{}

	for len(q) > 0 {
		p := q[0]
		q = q[1:]

		if cutoff > 0 && p.steps > cutoff {
			continue
		}

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
			if next.row < 0 || next.row >= len(grid) ||
				next.col < 0 || next.col >= len(grid[next.row]) {
				continue
			}

			if grid[next.row][next.col] == '#' {
				continue
			}

			if p.pos == c.start {
				q = append(q, node{c.end, cheat{}, p.steps + c.steps})
				continue
			}

			q = append(q, node{next, cheat{}, p.steps + 1})
		}
	}
	return -1
}

func contains(x Coord, l []Coord) bool {
	for _, c := range l {
		if c == x {
			return true
		}
	}

	return false
}
***/
