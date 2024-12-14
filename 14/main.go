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
	MAX_X = 101
	MAX_Y = 103
)

func init() {
	flag.BoolVar(&VERBOSE, "v", false, "print out extra debug info")

	flag.Parse()
	if flag.NArg() > 0 {
		FILE = flag.Arg(0)
		if FILE != "input" {
			MAX_X = 11
			MAX_Y = 7
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
	robots := []Robot{}
	for scanner.Scan() {
		line := scanner.Text()
		robots = append(robots, Parse(line))
	}

	print(robots)
	// part 1
	newR := make([]Robot, len(robots))
	for i, r := range robots {
		newR[i] = r.MoveN(100)
	}
	part1 := SafetyFactor(newR)

	part2 := 0
	minSF := part1
	for i := 0; i < MAX_X*MAX_Y; i++ {
		for j, r := range robots {
			newR[j] = r.MoveN(i)
		}

		sf := SafetyFactor(newR)
		if sf < minSF {
			debug("=========")
			debug(i, sf, minSF)
			print(newR)
			minSF = sf
		}

		if IsEasterEgg(newR) {
			part2 = i
			print(newR)
			break
		}
	}
	fmt.Println(part1, part2)
}

func SafetyFactor(robots []Robot) int {
	quadrants := CalculateQuadrants(robots)
	sf := 1
	for _, q := range quadrants {
		sf *= q
	}
	return sf
}

func print(robots []Robot) {
	if !VERBOSE {
		return
	}

	for i := 0; i < MAX_Y; i++ {
		for j := 0; j < MAX_X; j++ {
			count := 0
			for _, r := range robots {
				if r.p.x == j && r.p.y == i {
					count++
				}
			}
			if count > 0 {
				fmt.Printf("%d", count)
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

type Coord struct {
	x int
	y int
}

func (c Coord) Add(o Coord) Coord {
	return Coord{c.x + o.x, c.y + o.y}
}

func (c Coord) Scale(n int) Coord {
	return Coord{c.x * n, c.y * n}
}

type Robot struct {
	p Coord
	v Coord
}

func (r Robot) MoveN(n int) Robot {
	next := r.p.Add(r.v.Scale(n))
	next.x = (next.x + n*MAX_X) % MAX_X
	next.y = (next.y + n*MAX_Y) % MAX_Y

	return Robot{next, r.v}
}

func Parse(line string) Robot {
	var p, v Coord
	if _, err := fmt.Sscanf(line, "p=%d,%d v=%d,%d", &p.x, &p.y, &v.x, &v.y); err != nil {
		panic(err)
	}
	return Robot{p, v}
}

func CalculateQuadrants(robots []Robot) []int {
	quadrants := make([]int, 4)

	midX := (MAX_X + 1) / 2
	midY := (MAX_Y + 1) / 2

	for _, r := range robots {
		if r.p.x == midX-1 || r.p.y == midY-1 {
			continue
		}

		q := (r.p.x/midX)*2 + (r.p.y / (midY))
		quadrants[q]++
	}
	return quadrants
}

func IsEasterEgg(robots []Robot) bool {
	// All unique locations
	// seen := map[Coord]struct{}{}
	// for _, r := range robots {
	// 	if _, ok := seen[r.p]; ok {
	// 		return false
	// 	}
	// 	seen[r.p] = struct{}{}
	// }

	// return true

	quads := CalculateQuadrants(robots)
	for i, q := range quads {
		imbalanced := true
		for j, q2 := range quads {
			if i == j {
				continue
			}

			if q < q2*2 {
				imbalanced = false
				break
			}

		}

		if imbalanced {
			return true
		}
	}
	return false
}
