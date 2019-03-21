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
	tiles map[string]*tile
}

func (g *Group) NewTile(name string) {
	tiles[name] = &tile{
		index: len(g.nameMap),
		name:  name,
	}
}

func (g *Group) EdgeSet(tile string, edge Edge) *Set {
	t, ok := g.tiles[tile]
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
	bits   []uint64
	length int
	group  *Group
}

func newSet(g *Group) *Set {
	return &Set{group: g}
}

func (s *Set) Add(tile string) {
}

func (s *Set) Remove(tile string) {
}

func (s *Set) Has(tile string) bool {
	t, ok := s.group.tiles[tile]
	if !ok {
		return false
	}
	i := t.index
	if s.length <= i {
		return false
	}
	return s.bits[i/64]&(1<<(i%64)) == 1
}

func (s *Set) Clear() {
	s.bits = nil
	s.length = 0
}

type SetGrid struct {
}
