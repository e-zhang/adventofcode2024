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
	grid, moves, robot := parse(scanner)
	grid2, robot2, boxes := expand(grid)

	print(grid)
	for _, m := range moves {
		robot = Move(robot, grid, m)
		// debug("Move", string(m), ":")
		// print(grid)
	}
	fmt.Println(GPS(grid, 'O'))

	print(grid2)
	for i, m := range moves {
		robot2 = Move2(robot2, grid2, m)
		debug(i, "Move", string(m))
		print(grid2)
		invariants2(grid2, boxes)
	}
	fmt.Println(GPS(grid2, '['))

}

func parse(scanner *bufio.Scanner) ([][]byte, string, Coord) {
	grid := [][]byte{}
	var r Coord
	parseMoves := false
	var moves string
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			parseMoves = true
			continue
		}

		if !parseMoves {
			for i, c := range line {
				if c == '@' {
					r = Coord{len(grid), i}
				}
			}
			grid = append(grid, []byte(line))
		} else {
			moves += line
		}
	}

	return grid, moves, r
}

func expand(grid [][]byte) ([][]byte, Coord, int) {
	grid2 := [][]byte{}
	var r Coord
	boxes := 0
	for _, row := range grid {
		row2 := []byte{}
		for _, col := range row {
			switch col {
			case '.', '#':
				row2 = append(row2, col)
				row2 = append(row2, col)
			case '@':
				r = Coord{len(grid2), len(row2)}
				row2 = append(row2, col)
				row2 = append(row2, '.')
			case 'O':
				row2 = append(row2, '[')
				row2 = append(row2, ']')
				boxes++
			default:
				panic(col)
			}
		}
		grid2 = append(grid2, row2)
	}

	return grid2, r, boxes
}

func invariants2(grid [][]byte, boxes int) {
	count := 0
	for _, row := range grid {
		for _, col := range row {
			if col == '[' {
				count++
			}
		}
	}

	if count != boxes {
		panic(count)
	}
}

func print(grid [][]byte) {
	if !VERBOSE {
		return
	}

	for _, r := range grid {
		for _, c := range r {
			fmt.Printf(string(c))
		}
		fmt.Println()
	}
	fmt.Println()
}

func GPS(grid [][]byte, char byte) int {
	sum := 0
	for r, row := range grid {
		for c, col := range row {
			if col == char {
				gps := r*100 + c
				sum += gps
			}
		}
	}
	return sum
}

type Coord struct {
	row int
	col int
}

func (c Coord) Add(o Coord) Coord {
	return Coord{c.row + o.row, c.col + o.col}
}

func (c Coord) Scale(n int) Coord {
	return Coord{c.row * n, c.col * n}
}

func Dir(move rune) Coord {
	switch move {
	case '<':
		return Coord{0, -1}
	case '>':
		return Coord{0, 1}
	case '^':
		return Coord{-1, 0}
	case 'v':
		return Coord{1, 0}
	default:
		panic(move)
	}
}

func Move(robot Coord, grid [][]byte, move rune) Coord {
	dir := Dir(move)
	next := robot.Add(dir)
	curr := next
	var boxes []Coord
	for grid[curr.row][curr.col] == 'O' {
		boxes = append(boxes, curr)
		curr = curr.Add(dir)
	}

	if grid[curr.row][curr.col] == '.' {
		for _, b := range boxes {
			n := b.Add(dir)
			grid[n.row][n.col] = 'O'
		}

		grid[next.row][next.col] = '@'
		grid[robot.row][robot.col] = '.'
		return next
	}
	return robot
}

func isBox(v byte) bool {
	return v == '[' || v == ']'
}

func LookForBoxes(coords []Coord, grid [][]byte) bool {
	isWall := false
	hasBox := false
	for _, c := range coords {
		v := grid[c.row][c.col]
		if v == '#' {
			isWall = true
		}
		if isBox(v) {
			hasBox = true
		}
	}

	return !isWall && hasBox
}

func Move2(robot Coord, grid [][]byte, move rune) Coord {
	dir := Dir(move)
	next := robot.Add(dir)
	q := []Coord{next}
	var boxes []Coord
	for LookForBoxes(q, grid) {
		if dir.row == 0 {
			if len(q) != 1 {
				panic(q)
			}
			pair := q[0].Add(dir)
			boxes = append(boxes, q[0], pair)
			q = []Coord{pair.Add(dir)}
		} else {
			var n []Coord
			seen := map[Coord]struct{}{}
			for _, c := range q {
				if _, ok := seen[c]; ok {
					continue
				}

				v := grid[c.row][c.col]
				if !isBox(v) {
					continue
				}

				var pair Coord
				if v == '[' {
					pair = c.Add(Coord{0, 1})
				} else {
					pair = c.Add(Coord{0, -1})
				}

				seen[c] = struct{}{}
				seen[pair] = struct{}{}
				boxes = append(boxes, c, pair)
				n = append(n, c.Add(dir), pair.Add(dir))
			}
			q = n
		}
	}

	if !isEmpty(q, grid) {
		return robot
	}

	for i := len(boxes) - 1; i >= 0; i-- {
		b := boxes[i]
		n := b.Add(dir)
		grid[n.row][n.col] = grid[b.row][b.col]
		grid[b.row][b.col] = '.'
	}
	grid[next.row][next.col] = '@'
	grid[robot.row][robot.col] = '.'
	return next
}

func isEmpty(coords []Coord, grid [][]byte) bool {
	for _, c := range coords {
		v := grid[c.row][c.col]
		if v != '.' {
			return false
		}
	}
	return true
}
