#pragma once

#include <stdint.h>

void CmdHandler();
void CCmdHandler();

int HookedVFadeAlpha();
void HookedRDrawSequentialPoly(uintptr_t surf, int face);
void HookedRClear();
void HookedMemoryInit(uintptr_t buf, int size);

void HookedHUDRedraw(float time, int intermission);
void HookedHUDDrawTransparentTriangles();
int HookedHUDVidInit();
void HookedHUDReset();

void HookedPMInit(uintptr_t ppm);
void HookedPMPlayerMove(int server);
void __thiscall CHookedCGraphInitGraph(void *this);
void HookedCGraphInitGraph(uintptr_t this);
