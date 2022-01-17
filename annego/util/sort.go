package util

import (
	"sort"
	"reflect"
)

type IntSlice = sort.IntSlice
type Float64Slice = sort.Float64Slice
type StringSlice = sort.StringSlice

type ByteSlice []byte

func (p ByteSlice) Len() int           { return len(p) }
func (p ByteSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p ByteSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type Int8Slice []int8

func (p Int8Slice) Len() int           { return len(p) }
func (p Int8Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int8Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type Uint8Slice []uint8

func (p Uint8Slice) Len() int           { return len(p) }
func (p Uint8Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Uint8Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type Int16Slice []int16

func (p Int16Slice) Len() int           { return len(p) }
func (p Int16Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int16Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type Uint16Slice []uint16

func (p Uint16Slice) Len() int           { return len(p) }
func (p Uint16Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Uint16Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type Int32Slice []int32

func (p Int32Slice) Len() int           { return len(p) }
func (p Int32Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type Uint32Slice []uint32

func (p Uint32Slice) Len() int           { return len(p) }
func (p Uint32Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Uint32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type Int64Slice []int64

func (p Int64Slice) Len() int           { return len(p) }
func (p Int64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type Uint64Slice []uint64

func (p Uint64Slice) Len() int           { return len(p) }
func (p Uint64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Uint64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type Float32Slice []float32

func (p Float32Slice) Len() int           { return len(p) }
func (p Float32Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Float32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Sort any slice, s must be slice support
func Sort(s interface{}) {
	v := reflect.TypeOf(s)
	if v.Kind() != reflect.Slice {
		panic("sort must input slice")
	}
	switch v.Elem().Kind() {
	case reflect.Int:
		sort.Ints(s.([]int))
	case reflect.Int8:
		sort.Sort(Int8Slice(s.([]int8)))
	case reflect.Uint8:
		sort.Sort(Uint8Slice(s.([]uint8)))
	case reflect.Int16:
		sort.Sort(Int16Slice(s.([]int16)))
	case reflect.Uint16:
		sort.Sort(Uint16Slice(s.([]uint16)))
	case reflect.Int32:
		sort.Sort(Int32Slice(s.([]int32)))
	case reflect.Uint32:
		sort.Sort(Uint32Slice(s.([]uint32)))
	case reflect.Int64:
		sort.Sort(Int64Slice(s.([]int64)))
	case reflect.Uint64:
		sort.Sort(Uint64Slice(s.([]uint64)))
	case reflect.Float32:
		sort.Sort(Float32Slice(s.([]float32)))
	case reflect.Float64:
		sort.Sort(Float64Slice(s.([]float64)))
	case reflect.String:
		sort.Strings(s.([]string))
	default:
		panic("sort slice type unsupport: " + v.String())
	}
}