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

func (g *Group) NewTile(name string) {
	t = &tile{
		index: len(g.nameMap),
		name:  name,
	}
	g.tilesByName[name] = t
	g.tilesByIndex = append(g.tilesByIndex, t)
}

func (g *Group) EdgeSet(tile string, edge Edge) *Set {
	t, ok := g.tilesByName[tile]
	if !ok {
		return nil
	}
	e := t.edgeSets[edge]
	if e == nil {
		e = newSet(g)
		t.edgeSets[edge] = e
	}
	return e
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

func newSet(g *Group) *Set {
	return &Set{group: g}
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
	s.bits[i/64] |= uint64(1 << (i % 64))
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
	s.bits[i/64] &= ^uint64(1 << (i % 64))
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
	return s.bits[i/64]&uint64(1<<(i%64)) == 1
}

func (s *Set) Clear() {
	s.bits = nil
}

type SetGrid struct {
}
