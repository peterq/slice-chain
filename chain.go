package slice_chain

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"unsafe"
)

type anyTyp struct {
}

var anyType = reflect.TypeOf(anyTyp{})

type Collection struct {
	arr []interface{}
	typ reflect.Type
}

func Collect(src interface{}) Collection {
	rs := reflect.ValueOf(src)
	if rs.Type().Kind() != reflect.Slice {
		panic("src is not slice")
	}
	target := Collection{
		typ: rs.Type().Elem(),
	}
	for i := 0; i < rs.Len(); i++ {
		target.arr = append(target.arr, rs.Index(i).Interface())
	}
	return target
}

func (a Collection) Map(fn interface{}) Collection {
	rFn := checkFn(fn, []reflect.Type{a.typ}, []reflect.Type{reflect.TypeOf([]interface{}{}).Elem()})

	to := Collection{
		typ: rFn.Type().Out(0),
	}
	for _, val := range a.arr {
		ret := rFn.Call([]reflect.Value{reflect.ValueOf(val)})
		to.arr = append(to.arr, ret[0].Interface())
	}
	return to
}

func (a Collection) Filter(fn interface{}) Collection {
	rFn := checkFn(fn, []reflect.Type{a.typ}, []reflect.Type{reflect.TypeOf(false)})
	to := Collection{
		typ: a.typ,
	}
	for _, val := range a.arr {
		ret := rFn.Call([]reflect.Value{reflect.ValueOf(val)})
		if ret[0].Bool() {
			to.arr = append(to.arr, val)
		}
	}
	return to
}

func (a Collection) Sort(fn interface{}) Collection {
	rf := checkFn(fn, []reflect.Type{a.typ, a.typ}, []reflect.Type{reflect.TypeOf(false)})
	to := a.copy()
	sort.Slice(to.arr, func(i, j int) bool {
		l := to.arr[i]
		r := to.arr[j]
		ret := rf.Call([]reflect.Value{reflect.ValueOf(l), reflect.ValueOf(r)})
		return ret[0].Bool()
	})
	return to
}

func (a Collection) Find(startIndex int, fn interface{}, target interface{}) (index int) {
	index = -1
	rFn := checkFn(fn, []reflect.Type{a.typ}, []reflect.Type{reflect.TypeOf(false)})
	dest := reflect.Indirect(reflect.ValueOf(target))
	for i, v := range a.arr {
		if i < startIndex {
			continue
		}
		ret := rFn.Call([]reflect.Value{reflect.ValueOf(v)})
		if ret[0].Bool() {
			index = i
			dest.Set(reflect.ValueOf(v))
			break
		}
	}
	return
}

func (a Collection) IndexOf(v interface{}) int {
	for i, vv := range a.arr {
		if vv == v {
			return i
		}
	}
	return -1
}

func (a Collection) Reverse() Collection {
	to := a.copy()
	for i, j := 0, len(to.arr)-1; i < j; i, j = i+1, j-1 {
		to.arr[i], to.arr[j] = to.arr[j], to.arr[i]
	}
	return to
}

func (a Collection) Reduce(iv interface{}, to interface{}, fn interface{}) {
	rf := checkFn(fn, []reflect.Type{anyType, a.typ, reflect.TypeOf(0)},
		[]reflect.Type{reflect.TypeOf([]interface{}{}).Elem()})
	if !rf.Type().Out(0).AssignableTo(rf.Type().In(0)) {
		panic(fmt.Sprintf("return type must assignable to prev type, but %v cant assign to %v", rf.Type().Out(0), rf.Type().In(0)))
	}
	ri := reflect.ValueOf(iv)
	if !ri.Type().AssignableTo(rf.Type().In(0)) {
		panic("init value type must assignable to prev type")
	}
	dest := reflect.Indirect(reflect.ValueOf(to))
	dest.Set(reflect.ValueOf(iv))
	for idx, v := range a.arr {
		dest.Set(rf.Call([]reflect.Value{dest, reflect.ValueOf(v), reflect.ValueOf(idx)})[0])
	}
}

func (a Collection) SaveTo(ptr interface{}) {
	rp := reflect.ValueOf(ptr)
	if rp.Type().Kind() != reflect.Ptr {
		panic("ptr is not ptr")
	}
	if rp.Elem().Type().Elem() != a.typ && // type equal
		!(rp.Elem().Type().Elem().Kind() == reflect.Interface && // impl interface
			a.typ.Implements(rp.Elem().Type().Elem())) {
		log.Println(rp.Elem().Type(), a.typ)
		panic("slice element type invalid ")
	}
	newSlice := reflect.MakeSlice(rp.Type().Elem(), 0, 0)
	for _, v := range a.arr {
		newSlice = reflect.Append(newSlice, reflect.ValueOf(v))
	}
	dst := (*reflect.SliceHeader)(unsafe.Pointer(reflect.ValueOf(ptr).Pointer()))
	dst.Data = newSlice.Pointer()
	dst.Len = len(a.arr)
	dst.Cap = cap(a.arr)
}

func (a Collection) copy() Collection {
	to := Collection{
		typ: a.typ,
		arr: make([]interface{}, len(a.arr)),
	}
	copy(to.arr, a.arr)
	return to
}

func checkFn(fn interface{}, in []reflect.Type, out []reflect.Type) reflect.Value {
	rf := reflect.ValueOf(fn)
	if rf.Kind() != reflect.Func {
		panic("fn is not func")
	}
	if rf.Type().NumIn() != len(in) {
		panic(fmt.Sprintf("need %d in param(s), got %d", len(in), rf.Type().NumIn()))
	}
	if rf.Type().NumOut() != len(out) {
		panic(fmt.Sprintf("need %d out param(s), got %d", len(out), rf.Type().NumOut()))
	}
	for i, t := range in {
		if !t.AssignableTo(rf.Type().In(i)) && t != anyType {
			panic(fmt.Sprintf("in param with the index %d need %v, got %v", i, t, rf.Type().In(i)))
		}
	}
	for i, t := range out {
		if !rf.Type().Out(i).AssignableTo(t) && t != anyType {
			panic(fmt.Sprintf("out param with the index %d need %v, got %v", i, t, rf.Type().In(i)))
		}
	}
	return rf
}
