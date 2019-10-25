package slice_chain

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

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

type A struct {
	a int
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