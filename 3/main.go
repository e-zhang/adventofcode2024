package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
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
	var sum1, sum1r, sum2, sum2r int
	do, dor := true, true
	for scanner.Scan() {
		line := scanner.Text()

		var p1, p2 int

		p1, _ = parse(line, true, false)
		sum1 += p1

		p2, do = parse(line, do, true)
		sum2 += p2

		var r1, r2 int
		r1, _ = doRegex(line, true, false)
		sum1r += r1
		r2, dor = doRegex(line, dor, true)
		sum2r += r2
	}

	fmt.Println(sum1, sum2)
	fmt.Println("regex", sum1r, sum2r)
}

func doRegex(line string, do bool, part2 bool) (int, bool) {
	var re *regexp.Regexp
	if part2 {
		re = regexp.MustCompile(`mul\((?<first>\d{1,3}),(?<second>\d{1,3})\)|do\(\)|don't\(\)`)
	} else {
		re = regexp.MustCompile(`mul\((?<first>\d{1,3}),(?<second>\d{1,3})\)`)
	}

	matches := re.FindAllStringSubmatch(line, -1)

	var first, second, total int
	for _, m := range matches {
		if m[0] == "do()" {
			do = true
			continue
		}
		if m[0] == "don't()" {
			do = false
			continue
		}

		if !do {
			continue
		}

		if _, err := fmt.Sscanf(m[1], "%d", &first); err != nil {
			panic(err)
		}
		if _, err := fmt.Sscanf(m[2], "%d", &second); err != nil {
			panic(err)
		}

		total += first * second
	}

	return total, do
}
