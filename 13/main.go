package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
)

var (
	VERBOSE bool
	FILE    = "input"
)

const (
	OFFSET = 10000000000000
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

type Coord struct {
	x int
	y int
}

func (c Coord) Add(o Coord) Coord {
	return Coord{c.x + o.x, c.y + o.y}
}

type Machine struct {
	A     Coord
	B     Coord
	Prize Coord
}

func Parse(scanner *bufio.Scanner) Machine {
	line := scanner.Text()
	if line == "" {
		scanner.Scan()
		line = scanner.Text()
	}

	var a Coord
	if _, err := fmt.Sscanf(line, "Button A: X+%d, Y+%d", &a.x, &a.y); err != nil {
		panic(err)
	}
	scanner.Scan()
	line = scanner.Text()
	var b Coord
	if _, err := fmt.Sscanf(line, "Button B: X+%d, Y+%d", &b.x, &b.y); err != nil {
		panic(err)
	}

	scanner.Scan()
	line = scanner.Text()
	var p Coord
	if _, err := fmt.Sscanf(line, "Prize: X=%d, Y=%d", &p.x, &p.y); err != nil {
		panic(err)
	}

	return Machine{a, b, p}
}

type key struct {
	pos Coord
}

func main() {
	f, err := os.Open(FILE)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)
	// machines := []Machine{}
	tokens := 0
	for scanner.Scan() {
		m := Parse(scanner)
		debug(m)

		score := ScorePart2(m)
		if score > 0 {
			debug(m)
			tokens += score
		}
	}

	fmt.Println(tokens)
}

func ScorePart1(m Machine) int {
	min := -1
	for a := 0; a < 100; a++ {
		for b := 0; b < 100; b++ {
			if a*m.A.x+b*m.B.x == m.Prize.x &&
				a*m.A.y+b*m.B.y == m.Prize.y {
				score := a*3 + b
				if min < 0 || score < min {
					min = score
				}
			}
		}
	}

	return min
}

func ScorePart2(m Machine) int {
	m.Prize.x += OFFSET
	m.Prize.y += OFFSET

	//  system of equations
	// a * m.A.x + b * m.B.x = m.Prize.x
	// a * m.A.y + b * m.B.y = m.Prize.y

	// (m.Prize.y - b * m.B.y)/m.A.y = (m.Prize.x - b * m.B.x)/m.A.x

	// m.Prize.y / m.A.y - m.Prize.x / m.A.x = -b*m.B.x/m.A.x + b * m.B.y/m.A.y
	// m.Prize.y / m.A.y - m.Prize.x / m.A.x = b * (-m.B.x/m.A.x + m.B.y/m.A.y)

	bF := math.Round((float64(m.Prize.y)/float64(m.A.y) - float64(m.Prize.x)/float64(m.A.x)) / (-1*float64(m.B.x)/float64(m.A.x) + float64(m.B.y)/float64(m.A.y)))
	aF := math.Round((float64(m.Prize.x) - bF*float64(m.B.x)) / float64(m.A.x))
	a := int(aF)
	b := int(bF)

	if a*m.A.x+b*m.B.x == m.Prize.x &&
		a*m.A.y+b*m.B.y == m.Prize.y {
		return a*3 + b
	}

	return -1
}