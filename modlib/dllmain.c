#include "dllmain.h"

// Defined in Go.
extern void OnProcessAttach();
extern void OnProcessDetach();
extern uintptr_t GetLoadLibraryAAddr();
extern uintptr_t GetLoadLibraryWAddr();
extern void LoadLibraryACallback(LPCSTR);
extern void LoadLibraryWCallback(LPCWSTR);

static DWORD WINAPI OnAttachThread(LPVOID lpParam)
{
    OnProcessAttach();
    return 0;
}

static DWORD WINAPI OnDetachThread(LPVOID lpParam)
{
    OnProcessDetach();
    return 0;
}

BOOL WINAPI DllMain(HINSTANCE hinstDLL, DWORD fdwReason, LPVOID lpReserved)
{
    switch (fdwReason) {
    case DLL_PROCESS_ATTACH:
        CreateThread(NULL, 0, OnAttachThread, NULL, 0, NULL);
        break;
    case DLL_PROCESS_DETACH:
        CreateThread(NULL, 0, OnDetachThread, NULL, 0, NULL);
        break;
    }
    return TRUE;
}

// Have to write these entirely in C because we cannot define stdcall Go functions
HMODULE WINAPI HookedLoadLibraryA(LPCSTR lpLibFileName)
{
    HMODULE ret;
    ret = ((HMODULE WINAPI (*)(LPCSTR))GetLoadLibraryAAddr())(lpLibFileName);
    if (ret) {
        LoadLibraryACallback(lpLibFileName);
    }
    return ret;
}

HMODULE WINAPI HookedLoadLibraryW(LPCWSTR lpLibFileName)
{
    HMODULE ret;
    ret = ((HMODULE WINAPI (*)(LPCWSTR))GetLoadLibraryWAddr())(lpLibFileName);
    if (ret) {
        LoadLibraryWCallback(lpLibFileName);
    }
    return ret;
}
