#include "cdefs.h"

extern void CmdHandler();

void CCmdHandler()
{
    CmdHandler();
}

void __thiscall CHookedCGraphInitGraph(void *this)
{
    HookedCGraphInitGraph(this);
}
