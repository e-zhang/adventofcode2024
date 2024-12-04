package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)

	var grid []string

	for scanner.Scan() {
		line := scanner.Text()
		grid = append(grid, line)
	}

	var sum1, sum2 int
	for r, row := range grid {
		for c, l := range row {
			if l == 'X' {
				sum1 += part1(grid, r, c)
			}

			if l == 'A' {
				sum2 += part2(grid, r, c)
			}
		}
	}

	fmt.Println(sum1, sum2)
}

func part2(grid []string, r, c int) int {
	var corners []byte
	var ms, ss int
	for _, d := range []struct {
		x int
		y int
	}{
		{1, 1},   // DOWN RIGHT
		{1, -1},  // DOWN LEFT
		{-1, 1},  // UP RIGHT
		{-1, -1}, // UP LEFT
	} {
		x := r + d.x
		y := c + d.y

		if x < 0 || x >= len(grid) {
			break
		}
		if y < 0 || y >= len(grid[x]) {
			break
		}

		if grid[x][y] == 'M' {
			ms += 1
		} else if grid[x][y] == 'S' {
			ss += 1
		} else {
			break
		}

		corners = append(corners, grid[x][y])
	}

	if ms != 2 && ss != 2 {
		return 0
	}

	if len(corners) != 4 {
		return 0
	}

	if corners[0] == corners[3] || corners[1] == corners[2] {
		return 0
	}

	return 1
}

func part1(grid []string, r, c int) int {
	found := 0

	for _, d := range []struct {
		x int
		y int
	}{
		{0, 1},   // RIGHT
		{0, -1},  // LEFT
		{1, 0},   // DOWN
		{-1, 0},  // UP
		{1, 1},   // DOWN RIGHT
		{1, -1},  // DOWN LEFT
		{-1, 1},  // UP RIGHT
		{-1, -1}, // UP LEFT
	} {
		var word []byte
		for i := 0; i < 4; i++ {
			x := r + i*d.x
			y := c + i*d.y

			if x < 0 || x >= len(grid) {
				break
			}
			if y < 0 || y >= len(grid[x]) {
				break
			}

			word = append(word, grid[x][y])
		}
		if string(word) == "XMAS" {
			found++
		}
	}

	return found
}
