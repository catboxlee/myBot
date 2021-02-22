package helper

import (
	"math/rand"
	"reflect"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Min ...
func Min(n int, m int) int {
	if n > m {
		return m
	}
	return n
}

// Max ...
func Max(n int, m int) int {
	if n < m {
		return m
	}
	return n
}

// Abs ...
func Abs(n int) int {
	x := n >> 9
	return (n ^ x) - x
}

// InArray (Value, Slice)
func InArray(needle interface{}, haystack interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(haystack).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(haystack)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(needle, s.Index(i).Interface()) == true {
				exists = true
				index = i
				return
			}
		}
	}
	return
}

// Shuffle ...
func Shuffle(slice interface{}) {
	rv := reflect.ValueOf(slice)
	swap := reflect.Swapper(slice)
	length := rv.Len()
	for i := length - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		swap(i, j)
	}
}
