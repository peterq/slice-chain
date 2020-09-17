package slice_chain

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"
)

// string slice -> int slice -> string slice
func TestChain(t *testing.T) {
	var a = []string{"1", "2", "3", "4", "5"}

	Collect(a).
		Map(func(s string) int { i, _ := strconv.Atoi(s); return i }).
		Filter(func(i int) bool { return i%2 == 0 }).
		Map(func(i int) string { return fmt.Sprintf("%d * %d = %d", i, i, i*i) }).
		SaveTo(&a)

	if strings.Join(a, "; ") != "2 * 2 = 4; 4 * 4 = 16" {
		t.Fail()
	}
}

type intable interface {
	Int() int
}

type A struct {
	a int
}

func (a A) Int() int {
	return a.a
}

type B struct {
	c int
}

func TestStructSlice(t *testing.T) {

	var aa = []A{{1}, {2}, {3}}
	var bb []B

	Collect(aa).Map(func(a A) B { return B{a.a} }).SaveTo(&bb)

	if bb[2].c != 3 {
		t.Fail()
	}
}

type ia interface {
	a()
}
type a int

func (a) a() {
}

// convert real type slice to interface slice type
func TestSaveToInterface(t *testing.T) {
	var s []ia
	Collect([]a{1, 3, 4, 5}).SaveTo(&s)
	if len(s) != 4 {
		t.Fail()
	}
}

// sort struct
func TestCollection_Sort(t *testing.T) {
	var s = []A{{5}, {3}, {1}, {3}, {4}}
	Collect(s).Sort(func(l, r intable) bool { return l.Int() < r.Int() }).SaveTo(&s)
	last := s[0].a
	for _, v := range s {
		if v.a < last {
			t.Fail()
		}
		last = v.a
	}
}

// find element and index
func TestCollection_Find(t *testing.T) {
	var a intable
	var idx = Collect([]A{{5}, {3}, {1}, {3}, {4}}).
		Find(0, func(a A) bool { return a.a == 3 }, &a)
	if idx != 1 || a.Int() != 3 {
		t.Fail()
	}
}

func TestCollection_Reverse(t *testing.T) {
	var s = []string{"a", "b", "c"}
	Collect(s).Reverse().SaveTo(&s)
	if strings.Join(s, "") != "cba" {
		t.Fail()
	}
}

func TestCollection_Reduce(t *testing.T) {
	var r int
	Collect([]A{{1}, {2}, {3}, {4}, {5}}).
		Reduce(1, &r, func(prev int, cur intable, idx int) int { return prev * cur.Int() })

	if r != 120 {
		t.Fail()
	}
}

func TestCollection_Uniq(t *testing.T) {
	var a []A
	Collect([]A{{1}, {2}, {3}, {3}, {5}}).
		Uniq().SaveTo(&a)

	if fmt.Sprint(a) != "[{1} {2} {3} {5}]" {
		log.Println(a)
		t.FailNow()
	}
}

func TestCollectMap(t *testing.T) {
	var a = map[string]string{
		"a": "aa",
		"b": "bb",
		"c": "cc",
	}
	if len(CollectMapKeys(a).arr) != 3 || len(CollectMapValues(a).arr) != 3 {
		t.FailNow()
	}
}
