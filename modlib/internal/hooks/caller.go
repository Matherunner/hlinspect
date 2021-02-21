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

static inline int call_func_ints_3(void *fp, uintptr_t a, uintptr_t b, uintptr_t c)
{
	return ((int (*)(uintptr_t, uintptr_t, uintptr_t))fp)(a, b, c);
}

static inline int call_func_ints_4(void *fp, uintptr_t a, uintptr_t b, uintptr_t c, uintptr_t d)
{
	return ((int (*)(uintptr_t, uintptr_t, uintptr_t, uintptr_t))fp)(a, b, c, d);
}

static inline int call_func_float_int(void *fp, float a, uintptr_t b)
{
	return ((int (*)(float, uintptr_t))fp)(a, b);
}
*/
import "C"

func CallFuncInts0(address unsafe.Pointer) int {
	return int(C.call_func_ints_0(address))
}

func CallFuncInts1(address unsafe.Pointer, a uintptr) int {
	return int(C.call_func_ints_1(address, C.uint(a)))
}

func CallFuncInts3(address unsafe.Pointer, a, b, c uintptr) int {
	return int(C.call_func_ints_3(address, C.uint(a), C.uint(b), C.uint(c)))
}

func CallFuncInts4(address unsafe.Pointer, a, b, c, d uintptr) int {
	return int(C.call_func_ints_4(address, C.uint(a), C.uint(b), C.uint(c), C.uint(d)))
}

func CallFuncFloatInt(address unsafe.Pointer, a float32, b uintptr) int {
	return int(C.call_func_float_int(address, C.float(a), C.uint(b)))
}