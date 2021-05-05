#pragma once

#include <stdint.h>

void CCmdHandler();
void CmdHandler();

int HookedVFadeAlpha();
void HookedRDrawSequentialPoly(uintptr_t surf, int face);
void HookedRClear();
void HookedMemoryInit(uintptr_t buf, int size);
void HookedRDrawWorld();
void HookedModLoadBrushModel(uintptr_t model, uintptr_t buffer);
