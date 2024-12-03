package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	START = iota
	M
	U
	L
	OPEN
	COMMA
	D
	O
	N
	APOSTROPHE
	T
)

func part1(line string) int {
	var x int
	var first, second int
	cur := START

	var s int
	var total int
	for i, c := range line {
		switch c {
		case 'm':
			cur = M
			s = i
		case 'u':
			if cur == M {
				cur = U
			}
		case 'l':
			if cur == U {
				cur = L
			}
		case '(':
			if cur == L {
				cur = OPEN
				x = i
			}
		case ',':
			if cur == OPEN {
				cur = COMMA
				if _, err := fmt.Sscanf(line[x+1:i+1], "%d", &first); err != nil {
					panic(err)
				}
				x = i
			}
		case ')':
			if cur == COMMA {
				cur = START
				if _, err := fmt.Sscanf(line[x+1:i+1], "%d", &second); err != nil {
					panic(err)
				}

				total += first * second
				fmt.Println(string(line[s:i+1]), first, second, total)
			}
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if cur != OPEN && cur != COMMA {
				cur = START
			}
		default:
			cur = START
		}
	}

	return total
}

func part2(line string, do bool) (int, bool) {
	var start, paren, comma int

	var first, second int
	cur := START

	var total int
	for i, c := range line {
		switch c {
		case 'm', 'd':
			cur = M
			start = i
		case '(':
			if cur == M {
				cur = OPEN
				paren = i
			}
		case ',':
			if cur == OPEN {
				if _, err := fmt.Sscanf(line[paren+1:i], "%d", &first); err != nil {
					panic(err)
				}
				comma = i
			}
		case ')':
			if cur == OPEN {
				f := line[start:paren]
				fmt.Println("func", f, do)

				switch f {
				case "mul":
					if do && comma > paren {
						if _, err := fmt.Sscanf(line[comma+1:i], "%d", &second); err != nil {
							panic(err)
						}
						fmt.Println(string(line[start:i+1]), first, second, total)

						total += first * second
					}
				case "do":
					do = true
				case "don't":
					do = false
				}
				cur = START
			}
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			if cur != OPEN {
				cur = START
			}
		default:
			if cur != M {
				cur = START
			}
		}
	}

	return total, do
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)
	var sum1, sum2 int
	do := true
	for scanner.Scan() {
		line := scanner.Text()

		p1 := part1(line)
		sum1 += p1

		fmt.Println(">>>>>")
		p2, b := part2(line, do)
		do = b
		sum2 += p2

		fmt.Println("line", p1, p2)
	}

	fmt.Println(sum1, sum2)
}
