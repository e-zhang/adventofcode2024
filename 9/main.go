package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)

	// 1 long line
	scanner.Scan()
	line := scanner.Text()

	files := []int{}
	for i := 0; i < len(line); i += 2 {
		// fmt.Println(i)
		id := i / 2

		var blocks int
		if _, err := fmt.Sscanf(string(line[i]), "%d", &blocks); err != nil {
			panic(err)
		}

		for x := 0; x < blocks; x++ {
			files = append(files, id)
		}

		// odd number of spaces
		if i < len(line)-1 {
			var free int
			if _, err := fmt.Sscanf(string(line[i+1]), "%d", &free); err != nil {
				panic(err)
			}

			for x := 0; x < free; x++ {
				files = append(files, -1)
			}
		}
	}

	part1 := make([]int, len(files))
	copy(part1, files)

	part1 = moveBlocks(part1)
	files = moveFiles(files)

	// fmt.Println(files)
	fmt.Println(checksum(part1), checksum(files))
}

func checksum(files []int) int {
	sum := 0
	for i, x := range files {
		if x < 0 {
			continue
		}

		sum += i * x
	}
	return sum
}

// part1
func moveBlocks(files []int) []int {
	last := 0
	for i := len(files) - 1; i >= 0; i-- {
		val := files[i]
		space := -1
		for j := last; j < i; j++ {
			if files[j] == -1 {
				space = j
				break
			}
		}

		if space < 0 {
			break
		}

		last = space
		files[space] = val
		files[i] = -1
	}
	return files
}

// part2
func moveFiles(files []int) []int {
	for i := len(files) - 1; i >= 0; {
		fileEnd := i
		for ; fileEnd >= 0; fileEnd-- {
			if files[fileEnd] != files[i] {
				break
			}
		}

		f := files[fileEnd+1 : i+1]

		// look for empty from beginning
		for j := 0; j < i-len(f); {
			if files[j] >= 0 {
				j++
				continue
			}
			// found a free space
			freeEnd := j
			// look for enough free spaces for the whole file
			for ; freeEnd < j+len(f); freeEnd++ {
				if files[freeEnd] >= 0 {
					break
				}
			}

			// has enough free spaces
			if freeEnd == j+len(f) {
				copy(files[j:freeEnd], f)
				for x := range f {
					f[x] = -1
				}
				break
			}
			j = freeEnd
		}

		i = fileEnd
	}
	return files
}
