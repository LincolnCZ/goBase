package rbtree

import (
	"testing"
	"math/rand"
	"github.com/stretchr/testify/assert"
)

func testNewIntMap() Map {
	return NewMap(CompareInt)
}

func TestGetSetDelete(t *testing.T) {
	m := testNewIntMap()
	assert.EqualValues(t, m.Len(), 0)
	assert.False(t, m.Set(0, 10))
	assert.False(t, m.Set(1, 11))
	assert.False(t, m.Set(2, 12))
	assert.True(t, m.Set(0, 13))
	assert.EqualValues(t, m.Len(), 3)

	var val interface{}
	var ok bool

	val, ok = m.Get(-1)
	assert.EqualValues(t, val, nil)
	assert.EqualValues(t, ok, false)
	
	val, ok = m.Get(0)
	assert.EqualValues(t, val, 13)
	assert.EqualValues(t, ok, true)

	val, ok = m.Get(1)
	assert.EqualValues(t, val, 11)
	assert.EqualValues(t, ok, true)

	val, ok = m.Get(2)
	assert.EqualValues(t, val, 12)
	assert.EqualValues(t, ok, true)

	assert.EqualValues(t, m.DeleteWithKey(-1), false)
	assert.EqualValues(t, m.DeleteWithKey(1), true)
	_, ok = m.Get(1)
	assert.EqualValues(t, ok, false)
}

func TestOrder(t *testing.T) {
	// get random keys
	keys := make([]int, 0)
	for i := 0; i < 100; i++ {
		keys = append(keys, i)
	}
	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	m := testNewIntMap()
	for i, v := range keys {
		m.Set(v, i)
	}

	order := 0
	for iter := m.Min(); !iter.Equal(m.Limit()); iter = iter.Next() {
		idx := iter.Value().(int)
		assert.EqualValues(t, iter.Key(), order)
		assert.EqualValues(t, order, keys[idx])
		order++
	}

	order--
	for iter := m.Max(); !iter.Equal(m.NegativeLimit()); iter = iter.Prev() {
		idx := iter.Value().(int)
		assert.EqualValues(t, iter.Key(), order)
		assert.EqualValues(t, order, keys[idx])
		order--
	}
}

func TestFind(t *testing.T) {
	m := testNewIntMap()
	m.Set(0, nil)
	m.Set(2, nil)
	m.Set(3, nil)
	m.Set(7, nil)
	m.Set(9, nil)

	iter := m.Find(3)
	assert.EqualValues(t, iter.Key(), 3)
	assert.EqualValues(t, m.Find(1), m.Limit())

	iter = m.FindGE(3)
	assert.EqualValues(t, iter.Key(), 3)
	iter = m.FindGE(4)
	assert.EqualValues(t, iter.Key(), 7)
	assert.EqualValues(t, m.FindGE(10), m.Limit())

	iter = m.FindLE(3)
	assert.EqualValues(t, iter.Key(), 3)
	iter = m.FindLE(4)
	assert.EqualValues(t, iter.Key(), 3)
	iter = m.FindLE(10)
	assert.EqualValues(t, iter.Key(), 9)
	assert.EqualValues(t, m.FindLE(-1), m.NegativeLimit())
}