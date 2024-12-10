package main

import (
	"bufio"
	"fmt"
	"os"
)

type Coord struct {
	row int
	col int
}

func (c Coord) Add(o Coord) Coord {
	return Coord{c.row + o.row, c.col + o.col}
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)

	topMap := [][]int{}
	trailheads := []Coord{}

	for scanner.Scan() {
		line := scanner.Text()
		row := make([]int, len(line))
		for i, c := range line {
			if c == '.' {
				row[i] = -1
				continue
			}
			var h int
			if _, err := fmt.Sscanf(string(c), "%d", &h); err != nil {
				panic(err)
			}
			row[i] = h
			if h == 0 {
				trailheads = append(trailheads, Coord{len(topMap), i})
			}
		}

		topMap = append(topMap, row)
	}

	part1 := 0
	part2 := 0
	for _, t := range trailheads {
		s := score(topMap, t)
		r := rating(topMap, t)
		fmt.Println(t, s, r)
		part1 += s
		part2 += r
	}
	fmt.Println(part1, part2)
}

func score(topMap [][]int, pos Coord) int {
	return trails(topMap, pos, map[Coord]struct{}{})
}

func rating(topMap [][]int, pos Coord) int {
	return trails(topMap, pos, nil)
}

func trails(topMap [][]int, pos Coord, visit map[Coord]struct{}) int {
	h := topMap[pos.row][pos.col]
	if h == 9 {
		if visit != nil {
			if _, ok := visit[pos]; ok {
				return 0
			}

			visit[pos] = struct{}{}
		}
		return 1
	}

	s := 0
	for _, d := range []Coord{
		{0, -1}, // UP
		{0, 1},  //DOWN
		{-1, 0}, // LEFT
		{1, 0},  // RIGHT
	} {
		n := pos.Add(d)
		if n.row < 0 || n.row >= len(topMap) {
			continue
		}

		if n.col < 0 || n.col >= len(topMap[n.row]) {
			continue
		}

		hn := topMap[n.row][n.col]
		if hn == h+1 {
			s += trails(topMap, n, visit)
		}
	}

	return s
}
