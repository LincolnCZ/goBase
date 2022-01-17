package rbtree

type Pair struct {
	key Item
	value Item
}

// Map like map[interface{}]interface{}
// implement by Tree
type Map struct {
	tree *Tree
}

// NewMap Create a new empty Map
func NewMap(compare CompareFunc) Map {
	comparePair := func(a, b interface{}) int {
		pa := a.(Pair)
		pb := b.(Pair)
		return compare(pa.key, pb.key)
	}
	return Map{tree: NewTree(comparePair)}
}

// follow Map operation simple wrapper Tree

func (m Map) Len() int {
	return m.tree.Len()
}

func (m Map) Min() MapIterator {
	return MapIterator{Iterator: m.tree.Min()}
}

func (m Map) Max() MapIterator {
	return MapIterator{Iterator: m.tree.Max()}
}

func (m Map) Limit() MapIterator {
	return MapIterator{Iterator: m.tree.Limit()}
}

func (m Map) NegativeLimit() MapIterator {
	return MapIterator{Iterator: m.tree.NegativeLimit()}
}

func (m Map) Find(key Item) MapIterator {
	pair := Pair{key, nil}
	n, found := m.tree.findGE(pair)
	if !found {
		return MapIterator{Iterator{m.tree, nil}}
	}
	return MapIterator{Iterator{m.tree, n}}
}

func (m Map) FindGE(key Item) MapIterator {
	pair := Pair{key, nil}
	return MapIterator{Iterator: m.tree.FindGE(pair)}
}

func (m Map) FindLE(key Item) MapIterator {
	pair := Pair{key, nil}
	return MapIterator{Iterator: m.tree.FindLE(pair)}
}

// Get from map
// value: value find with key
// ok: ture if key found
func (m Map) Get(key Item) (value Item, ok bool) {
	pair := Pair{key, nil}
	n, found := m.tree.findGE(pair)
	if !found {
		return nil, false
	}
	pair = n.item.(Pair)
	return pair.value, true
}

// Set key and value, create new pair if not exist
// return true if key already exist
func (m Map) Set(key Item, value Item) bool {
	pair := Pair{key, value}
	n, found := m.tree.findGE(pair)
	if found {
		n.item = pair
	} else {
		m.tree.Insert(pair)
	}
	return found
}

func (m Map) DeleteWithKey(key Item) bool {
	pair := Pair{key, nil}
	n, found := m.tree.findGE(pair)
	if found {
		m.tree.doDelete(n)
		return true
	}
	return false
}

func (m Map) DeleteWithIterator(iter MapIterator) {
	m.tree.DeleteWithIterator(iter.Iterator)
}

func (m Map) Tree() *Tree {
	return m.tree
}

// MapIterator allows scanning map elements in sort order.
// implement by Iterator
type MapIterator struct {
	Iterator
}

func (iter MapIterator) Equal(iter2 MapIterator) bool {
	return iter.Iterator.Equal(iter2.Iterator)
}

func (iter MapIterator) Next() MapIterator {
	return MapIterator{Iterator: iter.Iterator.Next()}
}

func (iter MapIterator) Prev() MapIterator {
	return MapIterator{Iterator: iter.Iterator.Prev()}
}

func (iter MapIterator) Item() Pair {
	return iter.Iterator.Item().(Pair)
}

func (iter MapIterator) Key() Item {
	return iter.Item().key
}

func (iter MapIterator) Value() Item {
	return iter.Item().value
}