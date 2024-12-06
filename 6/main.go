package main

import (
	"bufio"
	"fmt"
	"os"
)

type Coord struct {
	x int
	y int
}

func (c Coord) Add(o Coord) Coord {
	return Coord{c.x + o.x, c.y + o.y}
}

func (c Coord) Rotate() Coord {
	// apply 90degree rotation matrix
	return Coord{-c.y, c.x}
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	grid := [][]byte{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		grid = append(grid, []byte(line))
	}

	start := findStart(grid)
	dir := Coord{0, -1} // UP

	part1(start, dir, grid)
	part2(start, dir, grid)
}

type Position struct {
	Coord
	dir Coord
}

func isLoop(c, dir Coord, grid [][]byte) bool {
	visit := map[Position]struct{}{}

	for {
		if _, ok := visit[Position{c, dir}]; ok {
			return true
		}
		visit[Position{c, dir}] = struct{}{}

		next := c.Add(dir)

		if next.y < 0 || next.y >= len(grid) {
			break
		}

		if next.x < 0 || next.x >= len(grid[next.y]) {
			break
		}

		if grid[next.y][next.x] == '#' {
			dir = dir.Rotate()
			continue
		}

		c = next
	}

	return false
}

func part2(start Coord, dir Coord, grid [][]byte) {
	loops := 0
	for y, r := range grid {
		for x := range r {
			orig := grid[y][x]
			grid[y][x] = '#'
			if isLoop(start, dir, grid) {
				loops++
			}
			grid[y][x] = orig
		}
	}

	fmt.Println(loops)
}

func part1(c Coord, dir Coord, grid [][]byte) {
	visit := map[Coord]struct{}{}
	visit[c] = struct{}{}

	for {
		next := c.Add(dir)

		if next.y < 0 || next.y >= len(grid) {
			break
		}

		if next.x < 0 || next.x >= len(grid[next.y]) {
			break
		}

		if grid[next.y][next.x] == '#' {
			dir = dir.Rotate()
			continue
		}

		visit[next] = struct{}{}
		c = next
	}

	fmt.Println(len(visit))
}

func printVisit(grid [][]byte, visit map[Coord]struct{}) {
	for c := range visit {
		grid[c.y][c.x] = 'X'
	}

	fmt.Println()
	for _, r := range grid {
		for _, c := range r {
			fmt.Printf("%s", string(c))
		}
		fmt.Println()
	}

	fmt.Println(len(visit))
}

func findStart(grid [][]byte) Coord {
	for i, r := range grid {
		for j, c := range r {
			if c == '^' {
				return Coord{j, i}
			}
		}
	}

	panic("no start")
}
