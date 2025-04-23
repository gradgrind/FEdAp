package main

import "fmt"

// #include <stdlib.h>
// #include "connector.h"
import "C"
import "unsafe"

var result string
var resultptr *C.char

// The //export allows C to call the Go func

//export GoCallback
func GoCallback(callback *C.char) *C.char {
    cb := C.GoString(callback)
	fmt.Printf("GoCallback got '%s'\n", cb)
	C.free(unsafe.Pointer(resultptr))
	result = HandleCallback(cb)
	resultptr = C.CString(result)
	return resultptr
}

func HandleCallback(data string) string {
	fmt.Printf("HandleCallback got '%s'\n", data)
	return "HandleCallback result"
}

func main() {
	fmt.Printf("Go says: calling C init ...\n")
	
	result = "Init message"
	resultptr = C.CString(result)
	
	C.init(resultptr)
	fmt.Printf("Go says: Finished\n")

	C.free(unsafe.Pointer(resultptr))
}
