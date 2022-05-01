#pragma once

#define WIN32_LEAN_AND_MEAN
#include <windows.h>

HMODULE WINAPI HookedLoadLibraryA(LPCSTR);
HMODULE WINAPI HookedLoadLibraryW(LPCWSTR);
