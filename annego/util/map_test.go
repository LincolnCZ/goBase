package util

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMapKey(t *testing.T) {
	m := make(map[int]string)
	m[0] = "0"
	m[1] = "1"
	m[2] = "2"
	v := GetMapKey(m).([]int)
	assert.NotNil(t, v)
	sort.Ints(v)
	assert.Equal(t, v, []int{0, 1, 2})
}

func TestMapSortedRange(t *testing.T) {
	m := make(map[int16]int32)
	for i := 0; i < 10; i++ {
		m[int16(i)] = int32(i)
	}
	index := 0
	MapSortedRange(m, func(k, v interface{}) {
		key := k.(int16)
		value := v.(int32)
		assert.True(t, index == int(key))
		assert.True(t, index == int(value))
		index++
	})
}

func TestMapDiff(t *testing.T) {
	m1 := make(map[int]int)
	m1[0] = 0
	m1[1] = 1
	m1[2] = 2
	m2 := make(map[int]string)

	r := MapDiff(m1, m2).([]int)
	sort.Ints(r)
	assert.Equal(t, r, []int{0, 1, 2})

	m2[0] = "1"
	r = MapDiff(m1, m2).([]int)
	sort.Ints(r)
	assert.Equal(t, r, []int{1, 2})
}
