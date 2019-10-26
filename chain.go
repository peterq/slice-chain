package slice_chain

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"unsafe"
)

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
	rFn := reflect.ValueOf(fn)
	if rFn.Type().Kind() != reflect.Func {
		panic("fn is not func")
	}
	if rFn.Type().NumIn() != 1 || rFn.Type().NumOut() != 1 {
		panic("fn in out param number must be one")
	}
	if rFn.Type().In(0) != a.typ {
		panic("fn in param type invalid")
	}
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
	rFn := reflect.ValueOf(fn)
	if rFn.Type().Kind() != reflect.Func {
		panic("fn is not func")
	}
	if rFn.Type().NumIn() != 1 || rFn.Type().NumOut() != 1 {
		panic("fn in out param number must be one")
	}
	if rFn.Type().In(0) != a.typ {
		panic("fn in param type invalid")
	}
	if rFn.Type().Out(0).Kind() != reflect.Bool {
		panic("fn out param type is not bool")
	}
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
	sort.Slice(a.arr, func(i, j int) bool {
		l := a.arr[i]
		r := a.arr[j]
		ret := rf.Call([]reflect.Value{reflect.ValueOf(l), reflect.ValueOf(r)})
		return ret[0].Bool()
	})
	return a
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
		if !t.AssignableTo(rf.Type().In(i)) {
			panic(fmt.Sprintf("in param with the index %d need %v, got %v", i, t, rf.Type().In(i)))
		}
	}
	for i, t := range out {
		if !rf.Type().Out(i).AssignableTo(t) {
			panic(fmt.Sprintf("out param with the index %d need %v, got %v", i, t, rf.Type().In(i)))
		}
	}
	return rf
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
