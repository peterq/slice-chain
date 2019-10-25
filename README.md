> chain style operation for golang's slice just like javascript array<br>
像 JavaScript 操作数组一样链式操作 golang 的slice


run:
```go
package main

import (
	"fmt"
	"github.com/peterq/slice-chain"
	"strconv"
	"strings"
)

func main() {
	var a = []string{"1", "2", "3", "4", "5"}

	slice_chain.Collect(a).
		Map(func(s string) int { i, _ := strconv.Atoi(s); return i }).
		Filter(func(i int) bool { return i%2 == 0 }).
		Map(func(i int) string { return fmt.Sprintf("%d * %d = %d", i, i, i*i) }).
		SaveTo(&a)

	println(strings.Join(a, "\n"))
}
```
got output: 
```
2 * 2 = 4
4 * 4 = 16
```

issue and pr is welcomed ;)

love golang, love china. may the world peace :)
