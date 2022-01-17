package queue

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func ensureEmpty(t *testing.T, q *Queue) {
	assert.Equal(t, q.Len(), 0)
	assert.Nil(t, q.Front())
	assert.Nil(t, q.Back())
	assert.Nil(t, q.PopFront())
	assert.Nil(t, q.PopBack())
}

func ensureEqual(t *testing.T, q *Queue, l []int) {
	assert.Equal(t, q.Len(), len(l))
	for i := 0; i < len(l); i++ {
		assert.Equal(t, q.Get(i), l[i])
	}
	assert.Nil(t, q.Get(len(l)))
}

func TestSingle(t *testing.T) {
	q := NewQueue()
	ensureEmpty(t, q)

	q.PushBack(0)
	assert.Equal(t, q.Len(), 1)
	assert.Equal(t, q.Front(), 0)
	assert.Equal(t, q.Back(), 0)

	q.PushFront(1)
	assert.Equal(t, q.Len(), 2)
	assert.Equal(t, q.Front(), 1)
	assert.Equal(t, q.Back(), 0)

	assert.Equal(t, q.PopFront(), 1)
	assert.Equal(t, q.PopBack(), 0)
	ensureEmpty(t, q)
}

func TestMult(t *testing.T) {
	q := NewQueue()
	l := make([]int, 0)

	for i := 0; i < 16; i++ {
		q.PushBack(i)
		l = append(l, i)
	}
	ensureEqual(t, q, l)

	for i := 0; i < 16; i++ {
		q.PushFront(i)
		l = append([]int{i}, l...)
	}
	ensureEqual(t, q, l)

	for i := 0; i < 8; i++ {
		q.PopFront()
		q.PopBack()
	}
	ensureEqual(t, q, l[8:24])

	for i := 0; i < 8; i++ {
		q.PopBack()
	}
	ensureEqual(t, q, l[8:16])
}