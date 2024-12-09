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

func (c Coord) Sub(o Coord) Coord {
	return Coord{
		c.row - o.row,
		c.col - o.col,
	}
}

func (c Coord) Scale(x int) Coord {
	return Coord{
		c.row * x,
		c.col * x,
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)
	freqs := map[rune][]Coord{}

	grid := []string{}
	for scanner.Scan() {
		line := scanner.Text()

		for i, c := range line {
			if c != '.' {
				if vals, ok := freqs[c]; ok {
					freqs[c] = append(vals, Coord{len(grid), i})
				} else {
					freqs[c] = []Coord{Coord{len(grid), i}}
				}
			}
		}
		grid = append(grid, line)
	}

	antinodes1 := map[Coord]struct{}{}
	antinodes2 := map[Coord]struct{}{}
	for _, locs := range freqs {
		for i, first := range locs {
			for _, second := range locs[i+1:] {
				d := first.Sub(second)

				antinodes2[first] = struct{}{}
				antinodes2[second] = struct{}{}
				i := 1
				for {
					a1 := second.Sub(d.Scale(i))
					if (a1.row < 0 || a1.row >= len(grid)) ||
						(a1.col < 0 || a1.col >= len(grid[a1.row])) {
						break
					}
					if i == 1 {
						antinodes1[a1] = struct{}{}
					}

					antinodes2[a1] = struct{}{}
					i++
				}

				i = 1
				for {
					a2 := first.Sub(d.Scale(-i))
					if (a2.row < 0 || a2.row >= len(grid)) ||
						(a2.col < 0 || a2.col >= len(grid[a2.row])) {
						break
					}
					if i == 1 {
						antinodes1[a2] = struct{}{}
					}

					antinodes2[a2] = struct{}{}
					i++
				}
			}
		}
	}

	print(grid, antinodes1)
	print(grid, antinodes2)
	fmt.Println(len(antinodes1), len(antinodes2))
}

func print(grid []string, nodes map[Coord]struct{}) {
	fmt.Println()
	for r, ln := range grid {
		for c, x := range ln {
			if _, ok := nodes[Coord{r, c}]; ok && x == '.' {
				fmt.Print("#")
			} else {
				fmt.Print(string(x))
			}
		}
		fmt.Println()
	}

}
