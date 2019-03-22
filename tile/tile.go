package tile

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
	gSize := s.group.Size()
	if gSize == 0 {
		return
	}
	s.bits = make([]uint64, (gSize-1)/64)
	fullInts := gSize / 64
	for i := 0; i < fullInts; i++ {
		s.bits[i] = uint64(1<<64 - 1)
	}
	if gSize%64 == 0 {
		return
	}
	s.bits[len(s.bits)-1] = (1<<uint(gSize%64) - 1)
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
	return s.bits[i/64]&uint64(1<<uint(i%64)) == 1
}

func (s *Set) Intersect(other *Set) {
}

func (s *Set) Union(other *Set) {
}

func (s *Set) Tiles() []string {
	tiles := []string{}
	for n := range s.group.tilesByName {
		if s.Has(n) {
			tiles = append(tiles, n)
		}
	}
	return tiles
}

func (s *Set) Clear() {
	s.bits = nil
}

func CollapseCell(grid [][]*Set, x, y int) {
}

func CollapseGrid(grid [][]*Set, startX, startY int) {
}

func GenerateGrid(width, height, randomSeed uint) (grid [][]*Set, resolved bool) {
	grid = make([][]*Set, width)

	for x := 0; x < width; x++ {
		newCol := make([]*Set, height)
		grid.cells[x] = newCol
		for y := 0; y < height; y++ {
			s := NewSet(group)
			s.AddAll()
			grid.cells[x][y] = s
		}
	}

}
