package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

var VERBOSE bool
var FILE = "input"

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
	wires := map[string]Wire{}
	initial := true
	z := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			initial = false
			continue
		}

		var w string
		var wire Wire
		if initial {
			tokens := strings.Split(line, ":")
			w = tokens[0]
			var v bool
			if _, err := fmt.Sscanf(tokens[1], " %t", &v); err != nil {
				panic(err)
			}
			wire = Value{v}
		} else {
			var in1, in2, op string
			if _, err := fmt.Sscanf(line, "%s %s %s -> %s", &in1, &op, &in2, &w); err != nil {
				panic(err)
			}
			wire = Gate{in1, in2, op}
		}
		wires[w] = wire

		if strings.HasPrefix(w, "z") {
			z = addBit(z, w)
		}
	}

	debug(z)

	zz := toNumber(z, wires)
	fmt.Println(zz)

	swaps := findSwaps(z, wires)
	sort.Strings(swaps)
	fmt.Println(strings.Join(swaps, ","))
}

func findSwaps(z []string, wires map[string]Wire) []string {
	swaps := []string{}

	// adder is
	// intermediates
	// a = x XOR y
	// b = a XOR cin
	// c = x AND y
	// cout = b OR c

	// z = a^cin
	var cin string
	for i, zb := range z {
		xb := fmt.Sprintf("x%02d", i)
		yb := fmt.Sprintf("y%02d", i)

		g := wires[zb].(Gate)
		if i == 0 {
			z := findGate(xb, yb, "XOR", wires)
			if z == "" {
				panic("bit0")
			}

			cin = findGate(xb, yb, "AND", wires)
			continue
		}

		a := findGate(xb, yb, "XOR", wires)
		s := findGate(cin, a, "XOR", wires)
		if s != zb {
			if s == "" {
				if g.in1 == cin {
					wires[g.in2], wires[a] = wires[a], wires[g.in2]
					swaps = append(swaps, a, g.in2)
					a = g.in2
				}
				if g.in2 == cin {
					wires[g.in1], wires[a] = wires[a], wires[g.in1]
					swaps = append(swaps, a, g.in1)
					a = g.in1
				}
			} else {
				wires[zb], wires[s] = wires[s], wires[zb]
				swaps = append(swaps, zb, s)
			}
		}
		b := findGate(a, cin, "AND", wires)
		c := findGate(xb, yb, "AND", wires)
		cout := findGate(c, b, "OR", wires)

		debug(a, b, c, cout, s)
		cin = cout
	}

	return swaps
}

func findGate(in1, in2, op string, wires map[string]Wire) string {
	for k, w := range wires {
		if strings.HasPrefix(k, "x") || strings.HasPrefix(k, "y") {
			continue
		}

		g := w.(Gate)
		if g.op != op {
			continue
		}

		if g.in1 != in1 && g.in1 != in2 {
			continue
		}

		if g.in2 != in2 && g.in2 != in1 {
			continue
		}

		return k
	}

	return ""
}

func addBit(x []string, w string) []string {
	var i int
	if _, err := fmt.Sscanf(w[1:], "%d", &i); err != nil {
		panic(err)
	}
	if i >= len(x) {
		x = append(x, make([]string, i-len(x)+1)...)
	}
	x[i] = w
	return x
}

func toNumber(x []string, wires map[string]Wire) int {
	num := 0
	for i, w := range x {
		if wires[w].Evaluate(wires) {
			num += 1 << i
		}
	}
	return num
}

type Wire interface {
	Evaluate(map[string]Wire) bool
}

type Value struct {
	v bool
}

func (v Value) Evaluate(_ map[string]Wire) bool {
	return v.v
}

type Gate struct {
	in1 string
	in2 string
	op  string
}

func (g Gate) Evaluate(wires map[string]Wire) bool {
	a := wires[g.in1].Evaluate(wires)
	b := wires[g.in2].Evaluate(wires)
	switch g.op {
	case "AND":
		return a && b
	case "OR":
		return a || b
	case "XOR":
		return a != b
	}

	panic(g.op)
}
