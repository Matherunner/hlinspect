package hooks

import "unsafe"

/*
#include <stdlib.h>

static inline int call_func_ints_0(void *fp)
{
	return ((int (*)())fp)();
}

static inline int call_func_ints_1(void *fp, uintptr_t a)
{
	return ((int (*)(uintptr_t))fp)(a);
}
*/
import "C"

func CallFuncInts0(address unsafe.Pointer) int {
	return int(C.call_func_ints_0(address))
}

func CallFuncInts1(address unsafe.Pointer, a uintptr) int {
	return int(C.call_func_ints_1(address, C.uint(a)))
}
