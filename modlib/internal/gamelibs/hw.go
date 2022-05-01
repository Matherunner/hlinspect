package gamelibs

import (
	"hlinspect/internal/engine"
	"hlinspect/internal/gamelibs/cdefs"
	"hlinspect/internal/hooks"
	"hlinspect/internal/logs"
	"unsafe"
)

var hwDLL *hooks.Module

func initHWRegistry(reg *APIRegistry) {
	reg.CvarRegisterVariable = hooks.MakeFunctionPattern("Cvar_RegisterVariable", nil, map[string]hooks.SearchPattern{
		HL8684: hooks.MustMakePattern("55 8B EC 83 EC 14 53 56 8B 75 08 57 8B 06 50 E8 ?? ?? ?? ?? 83 C4 04 85 C0 74 17 8B 0E 51 68"),
		HLNGHL: hooks.MustMakePattern("83 EC 14 53 56 8B 74 24 20 57 8B 06 50 E8 ?? ?? ?? ?? 83 C4 04 85 C0 74 17 8B 0E 51 68 ?? ?? ?? ?? E8 ?? ?? ?? ?? 83 C4 08 5F 5E 5B 83 C4 14 C3 8B 16 52 E8"),
	})
	reg.MemoryInit = hooks.MakeFunctionPattern("Memory_Init", nil, map[string]hooks.SearchPattern{
		HL8684: hooks.MustMakePattern("55 8B EC 8B 45 08 8B 4D 0C 56 BE 00 00 20 00 A3 ?? ?? ?? ?? 89 ?? ?? ?? ?? ?? C7 ?? ?? ?? ?? ?? ?? ?? ?? ?? C7 ?? ?? ?? ?? ?? ?? ?? ?? ?? E8 ?? ?? ?? ?? 68 ?? ?? ?? ?? E8"),
		HL4554: hooks.MustMakePattern("8B 44 24 04 8B 4C 24 08 56 BE 00 00 20 00 A3 ?? ?? ?? ?? 89 ?? ?? ?? ?? ?? C7 ?? ?? ?? ?? ?? ?? ?? ?? ?? C7 ?? ?? ?? ?? ?? ?? ?? ?? ?? E8 ?? ?? ?? ?? 68 ?? ?? ?? ?? E8"),
	})
	reg.CmdAddCommandWithFlags = hooks.MakeFunctionPattern("Cmd_AddCommandWithFlags", nil, map[string]hooks.SearchPattern{
		// Search for "Cmd_AddCommand: %s already defined as a var"
		HL8684: hooks.MustMakePattern("55 8B EC 56 57 8B 7D 08 57 E8 ?? ?? ?? ?? 8A 08 83 C4 04 84 C9 74 12 57 68 ?? ?? ?? ?? E8 ?? ?? ?? ?? 83 C4 08 5F 5E 5D C3 8B 35"),
		HLNGHL: hooks.MustMakePattern("56 57 8B 7C 24 0C 57 E8 ?? ?? ?? ?? 8A 08 83 C4 04 84 C9 74 11 57 68 ?? ?? ?? ?? E8 ?? ?? ?? ?? 83 C4 08 5F 5E C3 8B 35"),
	})
	reg.CmdArgv = hooks.MakeFunctionPattern("Cmd_Argv", nil, map[string]hooks.SearchPattern{
		// Search for "MISSING VALUE" to find Host_FullInfo_f
		// The first function called is Cmd_Argc, while the second function with one argument should be Cmd_Argv
		HL8684: hooks.MustMakePattern("55 8B EC 8D 45 08 50 FF 15 ?? ?? ?? ?? 8B 45 08 8B 0D ?? ?? ?? ?? 83 C4 04 3B C1 72 07 A1 ?? ?? ?? ?? 5D"),
		HL4554: hooks.MustMakePattern("8D 44 24 04 50 FF 15 ?? ?? ?? ?? 8B 44 24 08 8B 0D ?? ?? ?? ?? 83 C4 04 3B C1 72 06 A1 ?? ?? ?? ?? C3"),
	})
	reg.AngleVectors = hooks.MakeFunctionPattern("AngleVectors", nil, map[string]hooks.SearchPattern{
		HL8684: hooks.MustMakePattern("55 8B EC 83 EC 1C 8D 45 14 8D 4D 10 50 8D 55 0C 51 8D 45 08 52 50 FF 15 ?? ?? ?? ?? 8B 4D 08 83 C4 08"),
		HL4554: hooks.MustMakePattern("55 8B EC 83 E4 F8 83 EC 20 56 8D 45 14 57 8D 4D 10 50 8D 55 0C 51 8D 45 08 52 50 FF 15 ?? ?? ?? ?? 8B 4D 08 D9 41 04"),
		HLNGHL: hooks.MustMakePattern("55 8B EC 83 E4 F8 83 EC 20 8D 45 14 8D 4D 10 50 8D 55 0C 51 8D 45 08 52 50 FF 15 ?? ?? ?? ?? 8B 4D 08 83 C4 08"),
	})
	reg.PFTracelineDLL = hooks.MakeFunctionPattern("PF_traceline_DLL", nil, map[string]hooks.SearchPattern{
		HL8684: hooks.MustMakePattern("55 8B EC 8B 45 14 85 C0 75 05 A1 ?? ?? ?? ?? 8B 4D 0C 8B 55 08 56 50 8B 45 10 50 51 52 E8 ?? ?? ?? ?? D9 05"),
		HL4554: hooks.MustMakePattern("8B 44 24 10 85 C0 75 05 A1 ?? ?? ?? ?? 8B 4C 24 08 8B 54 24 04 56 50 8B 44 24 14 50 51 52 E8 ?? ?? ?? ?? D9 05"),
	})
	reg.BuildNumber = hooks.MakeFunctionPattern("build_number", nil, map[string]hooks.SearchPattern{
		HL8684: hooks.MustMakePattern("55 8B EC 83 EC 08 A1 ?? ?? ?? ?? 56 33 F6 85 C0 0F 85 9B 00 00 00 53 33 DB 8B 04 9D ?? ?? ?? ?? 8B 0D ?? ?? ?? ?? 6A 03 50 51 E8"),
		HL4554: hooks.MustMakePattern("A1 ?? ?? ?? ?? 83 EC 08 57 33 FF 85 C0 0F 85 A5 00 00 00 53 56 33 DB BE ?? ?? ?? ?? 8B 06 8B 0D"),
		HLNGHL: hooks.MustMakePattern("A1 ?? ?? ?? ?? 83 EC 08 56 33 F6 85 C0 0F 85 9F 00 00 00 53 33 DB 8B 04 9D ?? ?? ?? ?? 8B 0D"),
	})
	reg.VFadeAlpha = hooks.MakeFunctionPattern("V_FadeAlpha", nil, map[string]hooks.SearchPattern{
		HL8684: hooks.MustMakePattern("55 8B EC 83 EC 08 D9 05 ?? ?? ?? ?? DC 1D ?? ?? ?? ?? 8A 0D ?? ?? ?? ?? DF E0 F6 C4 05 7A 1C D9 05 ?? ?? ?? ?? DC 1D"),
		HL4554: hooks.MustMakePattern("D9 05 ?? ?? ?? ?? DC 1D ?? ?? ?? ?? 8A 0D ?? ?? ?? ?? 83 EC 08 DF E0 F6 C4"),
	})
	reg.DrawString = hooks.MakeFunctionPattern("Draw_String", nil, map[string]hooks.SearchPattern{
		// Search for "%i %i %i", there is a thunk that call this function. There are two functions that call the thunk.
		// The shorter one is Draw_String.
		HL8684: hooks.MustMakePattern("55 8B EC 56 57 E8 ?? ?? ?? ?? 8B 4D 0C 8B 75 08 50 8B 45 10 50 51 56 E8 ?? ?? ?? ?? 83 C4 10 8B F8 E8 ?? ?? ?? ?? 8D 04 37"),
		HL4554: hooks.MustMakePattern("56 57 E8 ?? ?? ?? ?? 8B 4C 24 10 8B 74 24 0C 50 8B 44 24 18 50 51 56 E8 ?? ?? ?? ?? 83 C4 10 8B F8 E8 ?? ?? ?? ?? 8D 04 37"),
	})
	reg.VGUI2DrawSetTextColorAlpha = hooks.MakeFunctionPattern("VGUI2_Draw_SetTextColorAlpha", nil, map[string]hooks.SearchPattern{
		HL8684: hooks.MustMakePattern("55 8B EC 8A 45 08 8A 4D 0C 8A 55 10 88 45 08 8A 45 14 88 4D 09 88 55 0A 88 45 0B 8B 4D 08 89"),
		HL4554: hooks.MustMakePattern("8A 44 24 04 8A 4C 24 08 8A 54 24 0C 88 44 24 04 8A 44 24 10 88 4C 24 05 88 54 24 06 88 44 24 07 8B 4C 24 04 89 0D"),
	})
	reg.HostAutoSaveF = hooks.MakeFunctionPattern("Host_AutoSave_f", nil, map[string]hooks.SearchPattern{
		HL8684: hooks.MustMakePattern("A1 ?? ?? ?? ?? B9 01 00 00 00 3B C1 0F 85 9F 00 00 00 A1 ?? ?? ?? ?? 85 C0 75 10 68 ?? ?? ?? ?? E8 ?? ?? ?? ?? 83 C4 04 33 C0 C3 39 0D"),
	})
	reg.HostNoclipF = hooks.MakeFunctionPattern("Host_Noclip_f", nil, map[string]hooks.SearchPattern{
		// Search for "noclip ON\n"
		HL8684: hooks.MustMakePattern("55 8B EC 83 EC 24 A1 ?? ?? ?? ?? BA 01 00 00 00 3B C2 75 09 E8 ?? ?? ?? ?? 8B E5 5D C3 D9 05 ?? ?? ?? ?? D8 1D"),
		HL4554: hooks.MustMakePattern("A1 ?? ?? ?? ?? BA 01 00 00 00 83 EC 24 3B C2 75 09 E8 ?? ?? ?? ?? 83 C4 24 C3 D9 05 ?? ?? ?? ?? D8 1D"),
		HLNGHL: hooks.MustMakePattern("A1 ?? ?? ?? ?? BA 01 00 00 00 83 EC 24 3B C2 75 08 83 C4 24 E9 ?? ?? ?? ?? D9 05 ?? ?? ?? ?? D8 1D"),
	})
	reg.TriGLRenderMode = hooks.MakeFunctionPattern("tri_GL_RenderMode", nil, map[string]hooks.SearchPattern{
		HL8684: hooks.MustMakePattern("55 8B EC 56 8B 75 08 83 FE 05 0F 87 ?? ?? ?? ?? FF 24 B5 ?? ?? ?? ?? 68 ?? ?? ?? ?? FF 15 ?? ?? ?? ?? 6A 01"),
		HL4554: hooks.MustMakePattern("56 8B 74 24 08 83 FE 05 0F 87 ?? ?? ?? ?? FF 24 B5 ?? ?? ?? ?? 68 ?? ?? ?? ?? FF 15 ?? ?? ?? ?? 6A 01"),
	})
	reg.TriGLBegin = hooks.MakeFunctionPattern("tri_GL_Begin", nil, map[string]hooks.SearchPattern{
		HL8684: hooks.MustMakePattern("55 8B EC E8 ?? ?? ?? ?? 8B 45 08 8B 0C 85 ?? ?? ?? ?? 51 FF 15 ?? ?? ?? ?? 5D C3"),
		HL4554: hooks.MustMakePattern("E8 ?? ?? ?? ?? 8B 44 24 04 8B 0C 85 ?? ?? ?? ?? 51 FF 15 ?? ?? ?? ?? C3"),
	})
	reg.TriGLEnd = hooks.MakeFunctionPattern("tri_GL_End", nil, map[string]hooks.SearchPattern{
		HL8684: hooks.MustMakePattern("FF 25 ?? ?? ?? ?? 90 90 90 90 90 90 90 90 90 90 55 8B EC 8B 45 0C"),
		HL4554: hooks.MustMakePattern("FF 25 ?? ?? ?? ?? 90 90 90 90 90 90 90 90 90 90 8B 44 24 08 8B 4C 24 04"),
	})
	reg.TriGLColor4f = hooks.MakeFunctionPattern("tri_GL_Color4f", nil, map[string]hooks.SearchPattern{
		HL8684: hooks.MustMakePattern("55 8B EC 51 83 3D ?? ?? ?? ?? 04 75 4A D9 45 14 D8 0D ?? ?? ?? ?? D9 5D FC D9 45 FC E8 ?? ?? ?? ?? D9 45 10"),
		HL4554: hooks.MustMakePattern("51 83 3D ?? ?? ?? ?? 04 75 50 D9 44 24 14 D8 0D ?? ?? ?? ?? D9 5C 24 00 D9 44 24 00 E8 ?? ?? ?? ?? D9 44 24 10"),
	})
	reg.TriGLCullFace = hooks.MakeFunctionPattern("tri_GL_CullFace", nil, map[string]hooks.SearchPattern{
		HL8684: hooks.MustMakePattern("55 8B EC 8B 45 08 83 E8 00 74 10 48 75 23 68 44 0B 00 00 FF 15 ?? ?? ?? ?? 5D C3 68 44 0B 00 00"),
		HL4554: hooks.MustMakePattern("8B 44 24 04 83 E8 00 74 0F 48 75 22 68 44 0B 00 00 FF 15 ?? ?? ?? ?? C3 68 44 0B 00 00"),
	})
	reg.TriGLVertex3fv = hooks.MakeFunctionPattern("tri_GL_Vertex3fv", nil, map[string]hooks.SearchPattern{
		HL8684: hooks.MustMakePattern("55 8B EC 8B 45 08 50 FF 15 ?? ?? ?? ?? 5D C3 90 55 8B EC 8B 45 10 8B 4D 0C 8B 55 08 50 51 52"),
		HL4554: hooks.MustMakePattern("8B 44 24 04 50 FF 15 ?? ?? ?? ?? C3 90 90 90 90 8B 44 24 0C 8B 4C 24 08 8B 54 24 04 50 51 52"),
	})
	reg.ScreenTransform = hooks.MakeFunctionPattern("ScreenTransform", nil, map[string]hooks.SearchPattern{
		HL8684: hooks.MustMakePattern("55 8B EC 51 8B 45 08 8B 4D 0C D9 05 ?? ?? ?? ?? D8 08 D9 05 ?? ?? ?? ?? D8 48 08 DE C1"),
		HL4554: hooks.MustMakePattern("51 8B 44 24 08 8B 4C 24 0C D9 05 ?? ?? ?? ?? D8 08 D9 05 ?? ?? ?? ?? D8 48 08 DE C1"),
	})
	reg.WorldTransform = hooks.MakeFunctionPattern("WorldTransform", nil, map[string]hooks.SearchPattern{
		// Most likely the function below ScreenTransform
		HL8684: hooks.MustMakePattern("55 8B EC 83 EC 08 8B 45 08 8B 4D 0C D9 05 ?? ?? ?? ?? D8 08 D9 05 ?? ?? ?? ?? D8 48"),
		HL4554: hooks.MustMakePattern("83 EC 08 8B 44 24 0C 8B 4C 24 10 D9 05 ?? ?? ?? ?? D8 08 D9 05 ?? ?? ?? ?? D8 48 08"),
	})
	reg.HudGetScreenInfo = hooks.MakeFunctionPattern("hudGetScreenInfo", nil, map[string]hooks.SearchPattern{
		// Search for "Half-Life %i/%s (hw build %d)". This function is Draw_ConsoleBackground
		// The function below it should be Draw_FillRGBA. Get the cross references to Draw_FillRGBA. One of them
		// is a global variable of enginefuncs. The next entry is hudGetScreenInfo.
		HL8684: hooks.MustMakePattern("55 8B EC 8D 45 08 50 FF 15 ?? ?? ?? ?? 8B 45 08 83 C4 04 85 C0 75 02 5D C3 81 38 14 02 00 00 74 04"),
		HL4554: hooks.MustMakePattern("8D 44 24 04 50 FF 15 ?? ?? ?? ?? 8B 44 24 08 83 C4 04 85 C0 75 01 C3 81 38 14 02 00 00 74 03"),
	})
	reg.RDrawSequentialPoly = hooks.MakeFunctionPattern("R_DrawSequentialPoly", nil, map[string]hooks.SearchPattern{
		HL8684: hooks.MustMakePattern("55 8B EC 51 A1 ?? ?? ?? ?? 53 56 57 83 B8 F8 02 00 00 01 75 63 E8 ?? ?? ?? ?? 68 03 03 00 00 68 02 03 00 00"),
		HL4554: hooks.MustMakePattern("A1 ?? ?? ?? ?? 53 55 56 8B 88 F8 02 00 00 BE 01 00 00 00 3B CE 57 75 ?? E8 ?? ?? ?? ?? 68 03 03 00 00 68 02 03 00 00"),
	})
	reg.RClear = hooks.MakeFunctionPattern("R_Clear", nil, map[string]hooks.SearchPattern{
		HL8684: hooks.MustMakePattern("8B 15 ?? ?? ?? ?? 33 C0 83 FA 01 0F 9F C0 50 E8 ?? ?? ?? ?? D9 05 ?? ?? ?? ?? DC 1D ?? ?? ?? ?? 83 C4 04 DF E0"),
		HLNGHL: hooks.MustMakePattern("D9 05 ?? ?? ?? ?? DC 1D ?? ?? ?? ?? DF E0 F6 C4 44 7B 34 D9 05 ?? ?? ?? ?? D8 1D"),
	})
	reg.PFCheckClientI = hooks.MakeFunctionPattern("PF_checkclient_I", nil, map[string]hooks.SearchPattern{
		// Search for "Spawned a NULL entity!", the referencing function is CreateNamedEntity
		// Find cross references, go to the global data, that data is g_engfuncsExportedToDlls
		// Go up 6 entries, we will end up at PF_checkclient_I
		HL8684: hooks.MustMakePattern("55 8B EC 83 EC 0C DD 05 ?? ?? ?? ?? DC 25 ?? ?? ?? ?? DC 1D ?? ?? ?? ?? DF E0 25 00 01 00 00 A1 ?? ?? ?? ?? 75 26"),
		HL4554: hooks.MustMakePattern("DD 05 ?? ?? ?? ?? DC 25 ?? ?? ?? ?? 83 EC 0C DC 1D ?? ?? ?? ?? DF E0 F6 C4 01 A1 ?? ?? ?? ?? 75 26"),
		HLNGHL: hooks.MustMakePattern("DD 05 ?? ?? ?? ?? DC 25 ?? ?? ?? ?? 83 EC 0C DC 1D ?? ?? ?? ?? DF E0 25 00 01 00 00 A1 ?? ?? ?? ?? 75 26"),
	})

	reg.CCmdHandler = cdefs.CDefs.CCmdHandler
}

func initHWDLL(base string) (err error) {
	if hwDLL != nil {
		return
	}

	hwDLL, err = hooks.NewModule(base)
	if err != nil {
		return
	}

	reg := Model.Registry()

	initHWRegistry(reg)

	items := map[*hooks.FunctionPattern]unsafe.Pointer{
		&reg.BuildNumber:                nil,
		&reg.CvarRegisterVariable:       nil,
		&reg.CmdAddCommandWithFlags:     nil,
		&reg.CmdArgv:                    nil,
		&reg.VFadeAlpha:                 cdefs.CDefs.HookedVFadeAlpha,
		&reg.DrawString:                 nil,
		&reg.VGUI2DrawSetTextColorAlpha: nil,
		&reg.HostAutoSaveF:              nil,
		&reg.HostNoclipF:                nil,
		&reg.PFTracelineDLL:             nil,
		&reg.TriGLRenderMode:            nil,
		&reg.TriGLBegin:                 nil,
		&reg.TriGLEnd:                   nil,
		&reg.TriGLColor4f:               nil,
		&reg.TriGLCullFace:              nil,
		&reg.TriGLVertex3fv:             nil,
		&reg.ScreenTransform:            nil,
		&reg.WorldTransform:             nil,
		&reg.HudGetScreenInfo:           nil,
		&reg.RClear:                     cdefs.CDefs.HookedRClear,
		&reg.RDrawSequentialPoly:        cdefs.CDefs.HookedRDrawSequentialPoly,
		&reg.MemoryInit:                 cdefs.CDefs.HookedMemoryInit,
		&reg.PFCheckClientI:             nil,
		&reg.AngleVectors:               nil,
	}

	errors := hooks.BatchFind(hwDLL, items)
	printBatchFindErrors(errors)

	switch reg.HostAutoSaveF.PatternKey() {
	case HL8684:
		ptr := *(*unsafe.Pointer)(unsafe.Pointer(uintptr(reg.HostAutoSaveF.Address()) + 19))
		engine.Engine.SetSV(ptr)
		logs.DLLLog.Debugf("Set SV address: %x", ptr)
	}

	switch reg.HostNoclipF.PatternKey() {
	case HL8684:
		ptr := *(*unsafe.Pointer)(unsafe.Pointer(uintptr(reg.HostNoclipF.Address()) + 31))
		ptr = unsafe.Pointer(uintptr(ptr) - 0x14)
		engine.Engine.SetGlobalVariables(ptr)
		logs.DLLLog.Debugf("Set GlobalVariables address: %x", ptr)
	case HL4554:
		ptr := *(*unsafe.Pointer)(unsafe.Pointer(uintptr(reg.HostNoclipF.Address()) + 28))
		ptr = unsafe.Pointer(uintptr(ptr) - 0x14)
		engine.Engine.SetGlobalVariables(ptr)
		logs.DLLLog.Debugf("Set GlobalVariables address: %x", ptr)
	case HLNGHL:
		ptr := *(*unsafe.Pointer)(unsafe.Pointer(uintptr(reg.HostNoclipF.Address()) + 27))
		ptr = unsafe.Pointer(uintptr(ptr) - 0x14)
		engine.Engine.SetGlobalVariables(ptr)
		logs.DLLLog.Debugf("Set GlobalVariables address: %x", ptr)
	}

	return nil
}
