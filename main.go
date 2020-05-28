package main

import (
	"bufio"
	"crypto/rand"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/oxplot/vense/tile"
)

var (
	sizeRegexp = regexp.MustCompile(`^\s*(\d+)\s*[xX]\s*(\d+)\s*$`)
	lineRegexp = regexp.MustCompile(`^\s*([^/=\s]+)\s*/\s*([trlb])\s*=\s*([^/=\s]+)/([trlb])\s*$`)

	gridSize   = &size{width: 10, height: 10}
	randomSeed = int64(0)
)

type size struct {
	width  int
	height int
}

func (s *size) String() string {
	return fmt.Sprintf("%dx%d", s.width, s.height)
}

func (s *size) Set(v string) error {
	m := sizeRegexp.FindStringSubmatch(v)
	if m == nil {
		return fmt.Errorf("grid size must be in wxh format")
	}
	w, _ := strconv.ParseUint(m[1], 10, 32)
	h, _ := strconv.ParseUint(m[2], 10, 32)
	if h == 0 || w == 0 {
		return fmt.Errorf("grid width and height must be > 0")
	}
	s.width = int(w)
	s.height = int(h)

	return nil
}

func parseInputLine(l string) (name1 string, e1 tile.Edge, name2 string, e2 tile.Edge, err error) {
	m := lineRegexp.FindStringSubmatch(l)
	if m == nil {
		err = fmt.Errorf("'%s' must be in form 'name/edge=name/edge'", l)
		return
	}
	name1 = m[1]
	name2 = m[3]
	switch m[2] {
	case "r":
		if m[4] != "l" {
			err = fmt.Errorf("right hand side of '%s' must have 'l' edge", l)
			return
		}
		e1 = tile.Right
		e2 = tile.Left
	case "l":
		if m[4] != "r" {
			err = fmt.Errorf("right hand side of '%s' must have 'r' edge", l)
			return
		}
		e1 = tile.Left
		e2 = tile.Right
	case "t":
		if m[4] != "b" {
			err = fmt.Errorf("right hand side of '%s' must have 'b' edge", l)
			return
		}
		e1 = tile.Top
		e2 = tile.Bottom
	case "b":
		if m[4] != "t" {
			err = fmt.Errorf("right hand side of '%s' must have 't' edge", l)
			return
		}
		e1 = tile.Bottom
		e2 = tile.Top
	}

	return
}

func run() error {

	tileGroup := tile.NewGroup()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		n1, e1, n2, e2, err := parseInputLine(scanner.Text())
		if err != nil {
			return err
		}
		tileGroup.Add(n1)
		tileGroup.Add(n2)
		tileGroup.EdgeSet(n1, e1).Add(n2)
		tileGroup.EdgeSet(n2, e2).Add(n1)
	}

	grid, ok := tile.GenerateGrid(gridSize.width, gridSize.height, tileGroup, randomSeed)
	if !ok {
		return fmt.Errorf("cannot collapse the grid onto a distinct state")
	}

	for y := 1; y < grid.Height(); y++ {
		fmt.Print(grid[0][y].Tiles()[0])
		for x := 1; x < grid.Width(); x++ {
			fmt.Printf(",%s", grid[x][y].Tiles()[0])
		}
		fmt.Print("\n")
	}

	return nil
}

func main() {

	flag.Var(gridSize, "size", "grid size in wxh - e.g. 12x14")
	flag.Int64Var(&randomSeed, "seed", 0, "random seed - if 0 or not specified, a random seed is used")
	flag.Parse()
	if randomSeed == 0 {
		var b [8]byte
		var err error
		if _, err = rand.Read(b[:]); err != nil {
			panic(err)
		}
		randomSeed = int64(binary.LittleEndian.Uint64(b[:]))
	}
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}
