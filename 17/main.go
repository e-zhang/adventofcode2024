package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	VERBOSE bool
	FILE    = "input"
)

const (
	ADV = iota
	BXL
	BST
	JNZ
	BXC
	OUT
	BDV
	CDV
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

func main() {
	f, err := os.Open(FILE)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	line := scanner.Text()
	var a, b, c int
	if _, err := fmt.Sscanf(line, "Register A: %d", &a); err != nil {
		panic(err)
	}
	scanner.Scan()
	line = scanner.Text()
	if _, err := fmt.Sscanf(line, "Register B: %d", &b); err != nil {
		panic(err)
	}
	scanner.Scan()
	line = scanner.Text()
	if _, err := fmt.Sscanf(line, "Register C: %d", &c); err != nil {
		panic(err)
	}

	scanner.Scan()
	scanner.Scan()
	line = scanner.Text()
	var prog string
	if _, err := fmt.Sscanf(line, "Program: %s", &prog); err != nil {
		panic(err)
	}

	var program []int
	for _, v := range strings.Split(prog, ",") {
		var i int
		if _, err := fmt.Sscanf(v, "%d", &i); err != nil {
			panic(err)
		}
		program = append(program, i)
	}

	debug(a, b, c, program)
	fmt.Println("part1", join(Run(program, a, b, c)))

	a = part2(program)
	fmt.Println("part2", a)
	if o := join(Run(program, a, b, c)); o != prog {
		panic(o)
	}
}

func cmp(x, y []int) bool {
	for i := 1; i <= len(x); i++ {
		if y[len(y)-i] != x[len(x)-i] {
			return false
		}
	}

	return true
}

func part2(program []int) int {
	a := 1
	var out []int
	for len(out) <= len(program) {
		out = Run(program, a, 0, 0)
		debug("run", a, out)
		if cmp(out, program) {
			if len(out) == len(program) {
				break
			}
			a <<= 3
			debug("equals", a)
			continue
		}
		a++
		debug("===", a, len(out))
	}

	return a
}

func combo(operand int, a, b, c int) int {
	switch operand {
	case 0, 1, 2, 3:
		return operand
	case 4:
		return a
	case 5:
		return b
	case 6:
		return c
	default:
		panic(operand)
	}
}

func Run(program []int, a, b, c int) []int {
	var out []int
	ptr := 0
	for ptr < len(program) {
		opcode, operand := program[ptr], program[ptr+1]
		jmp := false
		switch opcode {
		case ADV:
			a = a / (1 << combo(operand, a, b, c))
		case BXL:
			b = b ^ operand
		case BST:
			b = combo(operand, a, b, c) % 8
		case JNZ:
			if a != 0 {
				ptr = operand
				jmp = true
			}
		case BXC:
			b = b ^ c
		case OUT:
			out = append(out, combo(operand, a, b, c)%8)
		case BDV:
			b = a / (1 << combo(operand, a, b, c))
		case CDV:
			c = a / (1 << combo(operand, a, b, c))
		default:
			panic(opcode)
		}

		if !jmp {
			ptr += 2
		}
	}
	return out
}

func join(out []int) string {
	var s string
	for i, v := range out {
		if i != 0 {
			s += ","
		}
		s += fmt.Sprintf("%d", v)
	}
	return s
}
