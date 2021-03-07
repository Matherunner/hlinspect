#pragma once

#include <stdint.h>

void HookedPMInit(uintptr_t ppm);
void HookedPMPlayerMove(int server);
void __thiscall C_HookedCGraphInitGraph(void *this);
void HookedCGraphInitGraph(uintptr_t this);
