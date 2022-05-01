#include "cdefs.h"

void CCmdHandler()
{
    CmdHandler();
}

void __thiscall CHookedCGraphInitGraph(void *this)
{
    HookedCGraphInitGraph((uintptr_t)this);
}
