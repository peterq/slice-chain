> chain style operation for golang's slice just like javascript array

support operation:
 - Map:     convert slice of any type A to slice of any type B
 - Filter:  delete elements of a slice you do not need
 - Sort:    sort a slice of any type on your rule
 - Find:    find the element in the slice you need
 - Reverse: reverse the element
 
> 像 JavaScript 操作数组一样链式操作 golang 的slice

支持的操作:
 - Map:     从一种类型的切片转换成另一种类型的切片
 - Filter:  过滤切片中不符合条件的元素, 返回新的切片
 - Sort:    自定义规则对切片进行排序
 - Find:    找出切片中符合条件的第一个元素, 返回该元素和该元素的索引
 - Reverse: 颠倒切片元素先后顺序


eg: run:
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

> more usage you can find in the `chain_test.go` file


issue and pr are welcomed ;)

love golang, love china. may the world peace :)
