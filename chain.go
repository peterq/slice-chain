package slice_chain

import (
	"log"
	"reflect"
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
