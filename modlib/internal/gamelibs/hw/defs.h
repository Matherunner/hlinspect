#pragma once

int HookedVFadeAlpha();
void HookedRDrawSequentialPoly(uintptr_t surf, int face);
void HookedRClear();
void HookedMemoryInit(uintptr_t buf, int size);
