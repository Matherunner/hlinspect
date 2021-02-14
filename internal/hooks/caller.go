package hooks

import "unsafe"

/*
static inline int call_func_ints_0(void *fp)
{
	return ((int (*)())fp)();
}
*/
import "C"

func CallFuncInts0(address unsafe.Pointer) int {
	return int(C.call_func_ints_0(address))
}
