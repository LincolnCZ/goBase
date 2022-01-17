package util

import (
	"reflect"
)

// GetMapKey input map[key]value, get []key
func GetMapKey(m interface{}) interface{} {
	v := reflect.ValueOf(m)
	if v.Kind() != reflect.Map {
		panic("GetMapKey need input map")
	}

	ktype := v.Type().Key()
	r := reflect.MakeSlice(reflect.SliceOf(ktype), 0, v.Len())
	iter := v.MapRange()
	for iter.Next() {
		k := iter.Key()
		r = reflect.Append(r, k)
	}
	return r.Interface()
}

// MapSortedRange range map[K]V in order, call f(k, v)
func MapSortedRange(m interface{}, f func(k, v interface{})) {
	s := GetMapKey(m)
	Sort(s)
	v := reflect.ValueOf(m)
	r := reflect.ValueOf(s)

	for i := 0; i < r.Len(); i++ {
		k0 := r.Index(i)
		v0 := v.MapIndex(k0)
		f(k0.Interface(), v0.Interface())
	}
}

// MapDiff find key in m1 and not in m2, get []key
func MapDiff(m1, m2 interface{}) interface{} {
	v1 := reflect.ValueOf(m1)
	v2 := reflect.ValueOf(m2)
	if v1.Kind() != reflect.Map || v2.Kind() != reflect.Map {
		panic("MapDiff need input map")
	}
	if v1.Type().Key() != v2.Type().Key() {
		panic("MapDiff need input map of same key type")
	}

	ktype := v1.Type().Key()
	r := reflect.MakeSlice(reflect.SliceOf(ktype), 0, v1.Len())
	iter := v1.MapRange()
	for iter.Next() {
		k := iter.Key()
		if !v2.MapIndex(k).IsValid() {
			r = reflect.Append(r, k)
		}
	}
	return r.Interface()
}
