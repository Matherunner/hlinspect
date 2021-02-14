#include "dllmain.h"

extern void OnProcessAttach();
extern void OnProcessDetach();

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
