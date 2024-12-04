package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	START = iota
	M
	OPEN
	COMMA
)

func parse(line string, do bool, part2 bool) (int, bool) {
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
					if part2 {
						do = false
					}
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

		var p1, p2 int

		p1, _ = parse(line, true, false)
		sum1 += p1

		p2, do = parse(line, do, true)
		sum2 += p2
	}

	fmt.Println(sum1, sum2)
}
