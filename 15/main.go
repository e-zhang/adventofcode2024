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
	grid := [][]byte{}
	var robot Coord
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
					robot = Coord{len(grid), i}
				}
			}
			grid = append(grid, []byte(line))
		} else {
			moves += line
		}
	}
	debug(robot)

	grid2 := [][]byte{}
	var robot2 Coord
	boxes := 0
	for _, row := range grid {
		row2 := []byte{}
		for _, col := range row {
			switch col {
			case '.', '#':
				row2 = append(row2, col)
				row2 = append(row2, col)
			case '@':
				robot2 = Coord{len(grid2), len(row2)}
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
	debug(robot2, boxes, len(moves), GPS(grid2, '['))
	print(grid2)

	for i, m := range moves {
		robot2 = Move2(robot2, grid2, m)
		debug(i, "Move", string(m))
		print(grid2)
		if n := NumBoxes(grid2); n != boxes {
			panic(n)
		}
	}
	print(grid2)
	fmt.Println(GPS(grid2, '['))

	// for _, m := range moves {
	// 	robot = Move(robot, grid, m)
	// 	debug("Move", string(m), ":")
	// 	print(grid)
	// }

	// fmt.Println(GPS(grid, 'O'))
}

func NumBoxes(grid [][]byte) int {
	count := 0
	for _, row := range grid {
		for _, col := range row {
			if col == '[' {
				count++
			}
		}
	}

	return count
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

func Move2(robot Coord, grid [][]byte, move rune) Coord {
	dir := Dir(move)
	next := robot.Add(dir)
	curr := []Coord{next}
	var boxes []Coord
	for {
		isWall := false
		isBox := false
		for _, c := range curr {
			pos := grid[c.row][c.col]
			if pos == '#' {
				isWall = true
			}
			if pos == '[' || pos == ']' {
				isBox = true
			}
		}

		if isWall || !isBox {
			break
		}

		if dir.row == 0 {
			if len(curr) != 1 {
				panic(curr)
			}
			n := curr[0].Add(dir)
			boxes = append(boxes, curr[0], n)
			curr = []Coord{n.Add(dir)}
		} else {
			var nCurr []Coord
			seen := map[Coord]struct{}{}
			for _, c := range curr {
				if _, ok := seen[c]; ok {
					continue
				}

				pos := grid[c.row][c.col]
				if pos != '[' && pos != ']' {
					continue
				}

				var pair Coord
				if pos == '[' {
					pair = c.Add(Coord{0, 1})
				} else {
					pair = c.Add(Coord{0, -1})
				}

				seen[c] = struct{}{}
				seen[pair] = struct{}{}
				boxes = append(boxes, c, pair)
				nCurr = append(nCurr, c.Add(dir), pair.Add(dir))
			}
			curr = nCurr
		}
	}

	isEmpty := true
	for _, c := range curr {
		pos := grid[c.row][c.col]
		if pos != '.' {
			isEmpty = false
			break
		}
	}

	if isEmpty {
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
	return robot
}
