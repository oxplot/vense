package tile

import (
	"fmt"
	"math/bits"
	"math/rand"
	"strings"
)

type Edge int

const edgeCount = 4

const (
	Top Edge = iota
	Bottom
	Left
	Right
)

type Group struct {
	tilesByName  map[string]*tile
	tilesByIndex []*tile
}

func NewGroup() *Group {
	return &Group{
		tilesByName: map[string]*tile{},
	}
}

func (g *Group) Add(tileName string) {
	if g.Has(tileName) {
		return
	}
	t := &tile{
		index: g.Size(),
		name:  tileName,
	}
	g.tilesByName[tileName] = t
	g.tilesByIndex = append(g.tilesByIndex, t)
}

func (g *Group) Has(tileName string) bool {
	_, ok := g.tilesByName[tileName]
	return ok
}

func (g *Group) EdgeSet(tile string, edge Edge) *Set {
	t, ok := g.tilesByName[tile]
	if !ok {
		return nil
	}
	e := t.edgeSets[edge]
	if e == nil {
		e = NewSet(g)
		t.edgeSets[edge] = e
	}
	return e
}

func (g *Group) Size() int {
	return len(g.tilesByName)
}

type tile struct {
	index    int
	name     string
	edgeSets [edgeCount]*Set
}

type Set struct {
	bits  []uint64
	group *Group
}

func NewSet(group *Group) *Set {
	return &Set{group: group}
}

func (s *Set) Add(tile string) {
	t, ok := s.group.tilesByName[tile]
	if !ok {
		return
	}
	i := t.index
	if len(s.bits) <= (i / 64) {
		b := make([]uint64, i/64+1)
		copy(b, s.bits)
		s.bits = b
	}
	s.bits[i/64] |= uint64(1 << uint(i%64))
}

func (s *Set) AddAll() {
	for _, t := range s.group.tilesByIndex {
		s.Add(t.name)
	}
}

func (s *Set) Remove(tile string) {
	t, ok := s.group.tilesByName[tile]
	if !ok {
		return
	}
	i := t.index
	if len(s.bits) <= (i / 64) {
		return
	}
	s.bits[i/64] &= ^uint64(1 << uint(i%64))
}

func (s *Set) Has(tile string) bool {
	t, ok := s.group.tilesByName[tile]
	if !ok {
		return false
	}
	i := t.index
	if len(s.bits) <= (i / 64) {
		return false
	}
	return s.bits[i/64]&uint64(1<<uint(i%64)) > 0
}

func (s *Set) Intersect(other *Set) {
	minLen := len(s.bits)
	if len(other.bits) < minLen {
		minLen = len(other.bits)
	}
	for i := 0; i < minLen; i++ {
		s.bits[i] &= other.bits[i]
	}
	s.bits = s.bits[0:minLen]
}

func (s *Set) Union(other *Set) {
	minLen := len(s.bits)
	if len(other.bits) < minLen {
		minLen = len(other.bits)
	}
	for i := 0; i < minLen; i++ {
		s.bits[i] |= other.bits[i]
	}
	if len(other.bits) > len(s.bits) {
		s.bits = append(s.bits, other.bits[len(s.bits):]...)
	}
}

func (s *Set) Size() int {
	size := 0
	for _, b := range s.bits {
		size += bits.OnesCount64(b)
	}
	return size
}

func (s *Set) Tiles() []string {
	tiles := make([]string, 0, s.Size())
	for _, t := range s.group.tilesByIndex {
		i := t.index
		if len(s.bits) > (i/64) && s.bits[i/64]&uint64(1<<uint(i%64)) > 0 {
			tiles = append(tiles, t.name)
		}
	}
	return tiles
}

func (s *Set) tiles() []*tile {
	tiles := make([]*tile, 0, s.Size())
	for _, t := range s.group.tilesByIndex {
		i := t.index
		if len(s.bits) > (i/64) && s.bits[i/64]&uint64(1<<uint(i%64)) > 0 {
			tiles = append(tiles, t)
		}
	}
	return tiles
}

func (s *Set) Clear() {
	s.bits = nil
}

type Grid [][]*Set

func NewGrid(width, height int, group *Group) Grid {
	grid := Grid(make([][]*Set, width))

	for x := range grid {
		newCol := make([]*Set, height)
		grid[x] = newCol
		for y := range newCol {
			newCol[y] = NewSet(group)
		}
	}
	return grid
}

func (g Grid) Superposition() {
	for _, c := range g {
		for _, s := range c {
			s.AddAll()
		}
	}
}

func (g Grid) CollapseCell(x, y int) {
	tiles := g[x][y].tiles()
	group := g[0][0].group
	if x > 0 {
		s := NewSet(group)
		for _, t := range tiles {
			s.Union(t.edgeSets[Left])
		}
		g[x-1][y].Intersect(s)
	}
	if y > 0 {
		s := NewSet(group)
		for _, t := range tiles {
			s.Union(t.edgeSets[Top])
		}
		g[x][y-1].Intersect(s)
	}
	if x < len(g)-1 {
		s := NewSet(group)
		for _, t := range tiles {
			s.Union(t.edgeSets[Right])
		}
		g[x+1][y].Intersect(s)
	}
	if y < len(g[0])-1 {
		s := NewSet(group)
		for _, t := range tiles {
			s.Union(t.edgeSets[Bottom])
		}
		g[x][y+1].Intersect(s)
	}
}

func (g Grid) Collapse(startX, startY int) {
	for x := startX; x < len(g); x++ {
		for y := startY; y < len(g[x]); y++ {
			g.CollapseCell(x, y)
		}
		for y := startY - 1; y > 0; y-- {
			g.CollapseCell(x, y)
		}
	}
	for x := startX - 1; x > 0; x-- {
		for y := startY; y < len(g[x]); y++ {
			g.CollapseCell(x, y)
		}
		for y := startY - 1; y > 0; y-- {
			g.CollapseCell(x, y)
		}
	}
}

func GenerateGrid(width, height int, group *Group, randomSeed int64) (g Grid, ok bool) {
	rnd := rand.New(rand.NewSource(randomSeed))
	g = NewGrid(width, height, group)
	g.Superposition()

	firstX := rnd.Int() % width
	firstY := rnd.Int() % height
	firstCellTiles := g[firstX][firstY].Tiles()
	firstTile := firstCellTiles[rnd.Int()%len(firstCellTiles)]
	g[firstX][firstY].Clear()
	g[firstX][firstY].Add(firstTile)
	g.Collapse(firstX, firstY)

	//printGrid(g)

	lastX, lastY, lastSize := -1, -1, -1
	for {
		x, y, min, max := nextBestCell(g)
		if min == 1 && max == 1 {
			printGrid(g)
			return g, true
		}
		if min < 1 {
			return g, false
		}
		size := g[x][y].Size()
		if x == lastX && y == lastY && size == lastSize {
			return g, false
		}
		lastX, lastY, lastSize = x, y, size
		nextCellTiles := g[x][y].Tiles()
		picked := nextCellTiles[rnd.Int()%len(nextCellTiles)]
		g[x][y].Clear()
		g[x][y].Add(picked)
		g.Collapse(x, y)

		//printGrid(g)
	}
}

func nextBestCell(g Grid) (x, y, min, max int) {
	var bestX, bestY int
	bestMin := g[0][0].group.Size()
	min = bestMin
	for x, c := range g {
		for y, s := range c {
			ss := s.Size()
			if ss > 1 && ss < bestMin {
				bestMin = s.Size()
				bestX, bestY = x, y
			}
			if ss < min {
				min = ss
			}
			if ss > max {
				max = ss
			}
		}
	}
	return bestX, bestY, min, max
}

func printSet(s *Set) {
	fmt.Printf("[%s]", strings.Join(s.Tiles(), ","))
}

func printGrid(g Grid) {
	fmt.Println("---")
	for y := range g[0] {
		for x := range g {
			printSet(g[x][y])
			fmt.Print(" , ")
		}
		fmt.Print("\n")
	}
}

func PrintGroup(g *Group) {
	for t := range g.tilesByName {
		fmt.Printf("%s:\n", t)
		fmt.Printf("  left:   %s\n", strings.Join(g.EdgeSet(t, Left).Tiles(), ","))
		fmt.Printf("  right:  %s\n", strings.Join(g.EdgeSet(t, Right).Tiles(), ","))
		fmt.Printf("  top:    %s\n", strings.Join(g.EdgeSet(t, Top).Tiles(), ","))
		fmt.Printf("  bottom: %s\n", strings.Join(g.EdgeSet(t, Bottom).Tiles(), ","))
		fmt.Println()
	}
}
