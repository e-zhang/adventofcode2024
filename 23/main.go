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

func toString(v1, v2, v3 string) string {
	clique := []string{v1, v2, v3}
	sort.Strings(clique)
	return strings.Join(clique, ",")
}

func main() {
	f, err := os.Open(FILE)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)
	graph := map[string][]string{}
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, "-")
		v1, v2 := tokens[0], tokens[1]

		graph[v1] = append(graph[v1], v2)
		graph[v2] = append(graph[v2], v1)
	}

	fmt.Println(part1(graph), part1KB(graph))
	fmt.Println(part2(graph))
}

func part1(graph map[string][]string) int {
	cliques := map[string]struct{}{}
	seen := map[string]struct{}{}
	for v, edges := range graph {
		for i, e := range edges {
			if _, ok := seen[e]; ok {
				continue
			}
			for _, e2 := range edges[i+1:] {
				if _, ok := seen[e2]; ok {
					continue
				}
				if contains(e, graph[e2]) {
					cliques[toString(v, e, e2)] = struct{}{}
				}
				// seen[e2] = struct{}{}
			}
			// seen[e] = struct{}{}
		}
		seen[v] = struct{}{}
	}

	sets := []string{}
	for c := range cliques {
		sets = append(sets, c)
	}
	sort.Strings(sets)
	count := 0
	for _, c := range sets {
		debug(c)
		for _, v := range strings.Split(c, ",") {
			if v[0] == 't' {
				count++
				break
			}
		}
	}
	return count
}

func contains(x string, y []string) bool {
	for _, i := range y {
		if i == x {
			return true
		}
	}

	return false
}

// type Set map[string]struct{}
type Set []string

func (s Set) toString() string {
	// c := []string{}
	// for v := range s {
	// 	c = append(c, v)
	// }
	sort.Strings(s)
	return strings.Join(s, ",")
}

func (s Set) add(v string) Set {
	// n := map[string]struct{}{}
	// for k := range s {
	// 	n[k] = struct{}{}
	// }
	// n[v] = struct{}{}
	// return n
	n := make([]string, len(s))
	copy(n, s)
	return append(n, v)
}

func (s Set) intersect(o Set) Set {
	// n := map[string]struct{}{}
	// for k := range s {
	// 	if _, ok := o[k]; ok {
	// 		n[k] = struct{}{}
	// 	}
	// }
	// return n
	n := []string{}
	for _, v := range s {
		if contains(v, o) {
			n = append(n, v)
		}
	}
	return n
}

func (s Set) remove(v string) Set {
	// del(s, v)
	// return s

	n := []string{}
	for _, i := range s {
		if i != v {
			n = append(n, i)
		}
	}
	return n
}

func BronKerbosch(r, p, x Set, g map[string][]string) []string {
	if len(p) == 0 && len(x) == 0 {
		return []string{r.toString()}
	}

	cliques := []string{}
	pp := p
	for _, v := range p {
		c := BronKerbosch(r.add(v), pp.intersect(g[v]), x.intersect(g[v]), g)
		cliques = append(cliques, c...)
		pp = pp.remove(v)
		x = x.add(v)
	}

	return cliques
}

func part2(g map[string][]string) string {
	r := Set{}
	p := Set{}
	x := Set{}

	for v := range g {
		p = append(p, v)
	}

	cliques := BronKerbosch(r, p, x, g)

	var max string
	for _, c := range cliques {
		if len(c) > len(max) {
			max = c
		}
	}
	return max
}
