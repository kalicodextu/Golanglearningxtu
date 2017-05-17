//
//只有 tab 和 data 都为 nil 时，接口才等于 nil。
//
package main // tab = nil, data = nil
import (
	"fmt"
	"reflect"
	"unsafe"
)

var a interface{} = nil
var b interface{} = (*int)(nil) // tab 包含 *int 类型信息, data = nil
type iface struct {
	itab, data uintptr
}

func main() {
	ia := *(*iface)(unsafe.Pointer(&a))
	ib := *(*iface)(unsafe.Pointer(&b))
	fmt.Println(a == nil, ia)
	fmt.Println(b == nil, ib, reflect.ValueOf(b).IsNil())
}
