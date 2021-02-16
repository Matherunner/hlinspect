#pragma once

#define WIN32_LEAN_AND_MEAN
#include <windows.h>

void OnProcessAttach();
void OnProcessDetach();

uintptr_t GetLoadLibraryAAddr();
void LoadLibraryACallback(LPCSTR);
uintptr_t GetLoadLibraryWAddr();
void LoadLibraryWCallback(LPCWSTR);

HMODULE WINAPI HookedLoadLibraryA(LPCSTR);
HMODULE WINAPI HookedLoadLibraryW(LPCWSTR);
