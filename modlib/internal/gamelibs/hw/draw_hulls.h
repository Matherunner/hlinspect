#pragma once

#include "com_model.h"

typedef model_t *(*Mod_FindNameFunc)(int, const char *);

void hullLazyInit(Mod_FindNameFunc, int hullnum);
void hullDrawWorldHull(float alpha);
void hullClean();
void hullResetWorldModel(model_t *mod);