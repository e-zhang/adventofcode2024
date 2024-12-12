package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

var (
	NEIGHBORS = []Coord{
		{0, 1},
		{0, -1},
		{1, 0},
		{-1, 0},
	}
)

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	garden := []string{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		garden = append(garden, line)
	}

	total1 := 0
	total2 := 0
	regions := []Region{}
	for r, row := range garden {
		for c, col := range row {
			pos := Coord{r, c}

			// check to see if the pos has already been included in a region
			skip := false
			for _, region := range regions {
				if region.plant != col {
					continue
				}
				if _, ok := region.plots[pos]; ok {
					skip = true
					break
				}
			}

			if skip {
				continue
			}

			r := CreateRegion(col, pos, garden)
			a := r.Area()
			peri := r.Perimeter()
			price1 := a * peri
			sides := r.Sides()
			price2 := a * sides
			fmt.Println(string(r.plant), a, peri, sides)
			total1 += price1
			total2 += price2
			regions = append(regions, r)
		}
	}

	fmt.Println(total1, total2)
}

type Coord struct {
	row int
	col int
}

func (c Coord) Add(o Coord) Coord {
	return Coord{c.row + o.row, c.col + o.col}
}

type Region struct {
	plant rune
	plots map[Coord]struct{}
}

func (r Region) Area() int {
	return len(r.plots)
}

func (r Region) Perimeter() int {
	peri := 0
	for p := range r.plots {
		sides := 4
		for _, d := range NEIGHBORS {
			n := p.Add(d)
			if _, ok := r.plots[n]; ok {
				sides--
			}
		}
		peri += sides
	}

	return peri
}

type key struct {
	val int
	dir Coord
}

func (r Region) Sides() int {
	edges := map[key][]int{}
	for p := range r.plots {
		for _, d := range NEIGHBORS {
			n := p.Add(d)
			if _, ok := r.plots[n]; ok {
				continue
			}
			// looking at horizontal sides
			if d.row == 0 {
				k := key{n.col, d}
				if vals, ok := edges[k]; ok {
					edges[k] = append(vals, n.row)
				} else {
					edges[k] = []int{n.row}
				}
			}
			// looking at vertifcal sides
			if d.col == 0 {
				k := key{n.row, d}
				if vals, ok := edges[k]; ok {
					edges[k] = append(vals, n.col)
				} else {
					edges[k] = []int{n.col}
				}
			}
		}
	}
	sides := 0
	for _, vs := range edges {
		sort.Ints(vs)
		for i, v := range vs {
			// a gap indicates a new side
			if i == 0 || v-vs[i-1] > 1 {
				sides++
			}
		}
	}
	return sides
}

func CreateRegion(plant rune, pos Coord, garden []string) Region {
	q := []Coord{pos}
	plots := map[Coord]struct{}{}

	for len(q) != 0 {
		p := q[0]
		q = q[1:]

		if _, ok := plots[p]; ok {
			continue
		}
		plots[p] = struct{}{}

		for _, d := range NEIGHBORS {
			n := p.Add(d)

			if n.row < 0 || n.row >= len(garden) ||
				n.col < 0 || n.col >= len(garden[n.row]) {
				continue
			}

			if rune(garden[n.row][n.col]) != plant {
				continue
			}

			q = append(q, n)
		}
	}

	return Region{plant, plots}
}
