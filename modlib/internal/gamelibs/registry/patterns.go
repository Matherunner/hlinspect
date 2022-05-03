package registry

import (
	"hlinspect/internal/hooks"
)

// API holds the addresses to game DLL functions.
type API struct {
	// HW
	AngleVectors               hooks.FunctionPattern
	BuildNumber                hooks.FunctionPattern
	CmdAddCommandWithFlags     hooks.FunctionPattern
	CmdArgv                    hooks.FunctionPattern
	CvarRegisterVariable       hooks.FunctionPattern
	DrawString                 hooks.FunctionPattern
	HostAutoSaveF              hooks.FunctionPattern
	HostNoclipF                hooks.FunctionPattern
	HudGetScreenInfo           hooks.FunctionPattern
	MemoryInit                 hooks.FunctionPattern
	PFCheckClientI             hooks.FunctionPattern
	PFTracelineDLL             hooks.FunctionPattern
	RClear                     hooks.FunctionPattern
	RDrawSequentialPoly        hooks.FunctionPattern
	ScreenTransform            hooks.FunctionPattern
	TriGLBegin                 hooks.FunctionPattern
	TriGLColor4f               hooks.FunctionPattern
	TriGLCullFace              hooks.FunctionPattern
	TriGLEnd                   hooks.FunctionPattern
	TriGLRenderMode            hooks.FunctionPattern
	TriGLVertex3fv             hooks.FunctionPattern
	VFadeAlpha                 hooks.FunctionPattern
	VGUI2DrawSetTextColorAlpha hooks.FunctionPattern
	WorldTransform             hooks.FunctionPattern

	// CL
	HUDDrawTransparentTriangles hooks.FunctionPattern
	HUDRedraw                   hooks.FunctionPattern
	HUDReset                    hooks.FunctionPattern
	HUDVidInit                  hooks.FunctionPattern

	// HL
	CBaseMonsterChangeSchedule    hooks.FunctionPattern
	CBaseMonsterPBestSound        hooks.FunctionPattern
	CBaseMonsterRouteNew          hooks.FunctionPattern
	CGraphInitGraph               hooks.FunctionPattern
	CSoundEntActiveList           hooks.FunctionPattern
	CSoundEntSoundPointerForIndex hooks.FunctionPattern
	PMInit                        hooks.FunctionPattern
	PMPlayerMove                  hooks.FunctionPattern
	WorldGraph                    hooks.FunctionPattern
}

func NewAPI() *API {
	return &API{
		AngleVectors: hooks.MakeFunctionPattern("AngleVectors", nil, map[string]hooks.SearchPattern{
			HL8684: hooks.MustMakePattern("55 8B EC 83 EC 1C 8D 45 14 8D 4D 10 50 8D 55 0C 51 8D 45 08 52 50 FF 15 ?? ?? ?? ?? 8B 4D 08 83 C4 08"),
			HL4554: hooks.MustMakePattern("55 8B EC 83 E4 F8 83 EC 20 56 8D 45 14 57 8D 4D 10 50 8D 55 0C 51 8D 45 08 52 50 FF 15 ?? ?? ?? ?? 8B 4D 08 D9 41 04"),
			HLNGHL: hooks.MustMakePattern("55 8B EC 83 E4 F8 83 EC 20 8D 45 14 8D 4D 10 50 8D 55 0C 51 8D 45 08 52 50 FF 15 ?? ?? ?? ?? 8B 4D 08 83 C4 08"),
		}),
		BuildNumber: hooks.MakeFunctionPattern("build_number", nil, map[string]hooks.SearchPattern{
			HL8684: hooks.MustMakePattern("55 8B EC 83 EC 08 A1 ?? ?? ?? ?? 56 33 F6 85 C0 0F 85 9B 00 00 00 53 33 DB 8B 04 9D ?? ?? ?? ?? 8B 0D ?? ?? ?? ?? 6A 03 50 51 E8"),
			HL4554: hooks.MustMakePattern("A1 ?? ?? ?? ?? 83 EC 08 57 33 FF 85 C0 0F 85 A5 00 00 00 53 56 33 DB BE ?? ?? ?? ?? 8B 06 8B 0D"),
			HLNGHL: hooks.MustMakePattern("A1 ?? ?? ?? ?? 83 EC 08 56 33 F6 85 C0 0F 85 9F 00 00 00 53 33 DB 8B 04 9D ?? ?? ?? ?? 8B 0D"),
		}),
		CmdAddCommandWithFlags: hooks.MakeFunctionPattern("Cmd_AddCommandWithFlags", nil, map[string]hooks.SearchPattern{
			HL8684: hooks.MustMakePattern("55 8B EC 56 57 8B 7D 08 57 E8 ?? ?? ?? ?? 8A 08 83 C4 04 84 C9 74 12 57 68 ?? ?? ?? ?? E8 ?? ?? ?? ?? 83 C4 08 5F 5E 5D C3 8B 35"),
			HLNGHL: hooks.MustMakePattern("56 57 8B 7C 24 0C 57 E8 ?? ?? ?? ?? 8A 08 83 C4 04 84 C9 74 11 57 68 ?? ?? ?? ?? E8 ?? ?? ?? ?? 83 C4 08 5F 5E C3 8B 35"),
		}),
		CmdArgv: hooks.MakeFunctionPattern("Cmd_Argv", nil, map[string]hooks.SearchPattern{
			HL8684: hooks.MustMakePattern("55 8B EC 8D 45 08 50 FF 15 ?? ?? ?? ?? 8B 45 08 8B 0D ?? ?? ?? ?? 83 C4 04 3B C1 72 07 A1 ?? ?? ?? ?? 5D"),
			HL4554: hooks.MustMakePattern("8D 44 24 04 50 FF 15 ?? ?? ?? ?? 8B 44 24 08 8B 0D ?? ?? ?? ?? 83 C4 04 3B C1 72 06 A1 ?? ?? ?? ?? C3"),
		}),
		CvarRegisterVariable: hooks.MakeFunctionPattern("Cvar_RegisterVariable", nil, map[string]hooks.SearchPattern{
			HL8684: hooks.MustMakePattern("55 8B EC 83 EC 14 53 56 8B 75 08 57 8B 06 50 E8 ?? ?? ?? ?? 83 C4 04 85 C0 74 17 8B 0E 51 68"),
			HLNGHL: hooks.MustMakePattern("83 EC 14 53 56 8B 74 24 20 57 8B 06 50 E8 ?? ?? ?? ?? 83 C4 04 85 C0 74 17 8B 0E 51 68 ?? ?? ?? ?? E8 ?? ?? ?? ?? 83 C4 08 5F 5E 5B 83 C4 14 C3 8B 16 52 E8"),
		}),
		DrawString: hooks.MakeFunctionPattern("Draw_String", nil, map[string]hooks.SearchPattern{
			HL8684: hooks.MustMakePattern("55 8B EC 56 57 E8 ?? ?? ?? ?? 8B 4D 0C 8B 75 08 50 8B 45 10 50 51 56 E8 ?? ?? ?? ?? 83 C4 10 8B F8 E8 ?? ?? ?? ?? 8D 04 37"),
			HL4554: hooks.MustMakePattern("56 57 E8 ?? ?? ?? ?? 8B 4C 24 10 8B 74 24 0C 50 8B 44 24 18 50 51 56 E8 ?? ?? ?? ?? 83 C4 10 8B F8 E8 ?? ?? ?? ?? 8D 04 37"),
		}),
		HostAutoSaveF: hooks.MakeFunctionPattern("Host_AutoSave_f", nil, map[string]hooks.SearchPattern{
			HL8684: hooks.MustMakePattern("A1 ?? ?? ?? ?? B9 01 00 00 00 3B C1 0F 85 9F 00 00 00 A1 ?? ?? ?? ?? 85 C0 75 10 68 ?? ?? ?? ?? E8 ?? ?? ?? ?? 83 C4 04 33 C0 C3 39 0D"),
		}),
		HostNoclipF: hooks.MakeFunctionPattern("Host_Noclip_f", nil, map[string]hooks.SearchPattern{
			HL8684: hooks.MustMakePattern("55 8B EC 83 EC 24 A1 ?? ?? ?? ?? BA 01 00 00 00 3B C2 75 09 E8 ?? ?? ?? ?? 8B E5 5D C3 D9 05 ?? ?? ?? ?? D8 1D"),
			HL4554: hooks.MustMakePattern("A1 ?? ?? ?? ?? BA 01 00 00 00 83 EC 24 3B C2 75 09 E8 ?? ?? ?? ?? 83 C4 24 C3 D9 05 ?? ?? ?? ?? D8 1D"),
			HLNGHL: hooks.MustMakePattern("A1 ?? ?? ?? ?? BA 01 00 00 00 83 EC 24 3B C2 75 08 83 C4 24 E9 ?? ?? ?? ?? D9 05 ?? ?? ?? ?? D8 1D"),
		}),
		HudGetScreenInfo: hooks.MakeFunctionPattern("hudGetScreenInfo", nil, map[string]hooks.SearchPattern{
			HL8684: hooks.MustMakePattern("55 8B EC 8D 45 08 50 FF 15 ?? ?? ?? ?? 8B 45 08 83 C4 04 85 C0 75 02 5D C3 81 38 14 02 00 00 74 04"),
			HL4554: hooks.MustMakePattern("8D 44 24 04 50 FF 15 ?? ?? ?? ?? 8B 44 24 08 83 C4 04 85 C0 75 01 C3 81 38 14 02 00 00 74 03"),
		}),
		MemoryInit: hooks.MakeFunctionPattern("Memory_Init", nil, map[string]hooks.SearchPattern{
			HL8684: hooks.MustMakePattern("55 8B EC 8B 45 08 8B 4D 0C 56 BE 00 00 20 00 A3 ?? ?? ?? ?? 89 ?? ?? ?? ?? ?? C7 ?? ?? ?? ?? ?? ?? ?? ?? ?? C7 ?? ?? ?? ?? ?? ?? ?? ?? ?? E8 ?? ?? ?? ?? 68 ?? ?? ?? ?? E8"),
			HL4554: hooks.MustMakePattern("8B 44 24 04 8B 4C 24 08 56 BE 00 00 20 00 A3 ?? ?? ?? ?? 89 ?? ?? ?? ?? ?? C7 ?? ?? ?? ?? ?? ?? ?? ?? ?? C7 ?? ?? ?? ?? ?? ?? ?? ?? ?? E8 ?? ?? ?? ?? 68 ?? ?? ?? ?? E8"),
		}),
		PFCheckClientI: hooks.MakeFunctionPattern("PF_checkclient_I", nil, map[string]hooks.SearchPattern{
			HL8684: hooks.MustMakePattern("55 8B EC 83 EC 0C DD 05 ?? ?? ?? ?? DC 25 ?? ?? ?? ?? DC 1D ?? ?? ?? ?? DF E0 25 00 01 00 00 A1 ?? ?? ?? ?? 75 26"),
			HL4554: hooks.MustMakePattern("DD 05 ?? ?? ?? ?? DC 25 ?? ?? ?? ?? 83 EC 0C DC 1D ?? ?? ?? ?? DF E0 F6 C4 01 A1 ?? ?? ?? ?? 75 26"),
			HLNGHL: hooks.MustMakePattern("DD 05 ?? ?? ?? ?? DC 25 ?? ?? ?? ?? 83 EC 0C DC 1D ?? ?? ?? ?? DF E0 25 00 01 00 00 A1 ?? ?? ?? ?? 75 26"),
		}),
		PFTracelineDLL: hooks.MakeFunctionPattern("PF_traceline_DLL", nil, map[string]hooks.SearchPattern{
			HL8684: hooks.MustMakePattern("55 8B EC 8B 45 14 85 C0 75 05 A1 ?? ?? ?? ?? 8B 4D 0C 8B 55 08 56 50 8B 45 10 50 51 52 E8 ?? ?? ?? ?? D9 05"),
			HL4554: hooks.MustMakePattern("8B 44 24 10 85 C0 75 05 A1 ?? ?? ?? ?? 8B 4C 24 08 8B 54 24 04 56 50 8B 44 24 14 50 51 52 E8 ?? ?? ?? ?? D9 05"),
		}),
		RClear: hooks.MakeFunctionPattern("R_Clear", nil, map[string]hooks.SearchPattern{
			HL8684: hooks.MustMakePattern("8B 15 ?? ?? ?? ?? 33 C0 83 FA 01 0F 9F C0 50 E8 ?? ?? ?? ?? D9 05 ?? ?? ?? ?? DC 1D ?? ?? ?? ?? 83 C4 04 DF E0"),
			HLNGHL: hooks.MustMakePattern("D9 05 ?? ?? ?? ?? DC 1D ?? ?? ?? ?? DF E0 F6 C4 44 7B 34 D9 05 ?? ?? ?? ?? D8 1D"),
		}),
		RDrawSequentialPoly: hooks.MakeFunctionPattern("R_DrawSequentialPoly", nil, map[string]hooks.SearchPattern{
			HL8684: hooks.MustMakePattern("55 8B EC 51 A1 ?? ?? ?? ?? 53 56 57 83 B8 F8 02 00 00 01 75 63 E8 ?? ?? ?? ?? 68 03 03 00 00 68 02 03 00 00"),
			HL4554: hooks.MustMakePattern("A1 ?? ?? ?? ?? 53 55 56 8B 88 F8 02 00 00 BE 01 00 00 00 3B CE 57 75 ?? E8 ?? ?? ?? ?? 68 03 03 00 00 68 02 03 00 00"),
		}),
		ScreenTransform: hooks.MakeFunctionPattern("ScreenTransform", nil, map[string]hooks.SearchPattern{
			HL8684: hooks.MustMakePattern("55 8B EC 51 8B 45 08 8B 4D 0C D9 05 ?? ?? ?? ?? D8 08 D9 05 ?? ?? ?? ?? D8 48 08 DE C1"),
			HL4554: hooks.MustMakePattern("51 8B 44 24 08 8B 4C 24 0C D9 05 ?? ?? ?? ?? D8 08 D9 05 ?? ?? ?? ?? D8 48 08 DE C1"),
		}),
		TriGLBegin: hooks.MakeFunctionPattern("tri_GL_Begin", nil, map[string]hooks.SearchPattern{
			HL8684: hooks.MustMakePattern("55 8B EC E8 ?? ?? ?? ?? 8B 45 08 8B 0C 85 ?? ?? ?? ?? 51 FF 15 ?? ?? ?? ?? 5D C3"),
			HL4554: hooks.MustMakePattern("E8 ?? ?? ?? ?? 8B 44 24 04 8B 0C 85 ?? ?? ?? ?? 51 FF 15 ?? ?? ?? ?? C3"),
		}),
		TriGLColor4f: hooks.MakeFunctionPattern("tri_GL_Color4f", nil, map[string]hooks.SearchPattern{
			HL8684: hooks.MustMakePattern("55 8B EC 51 83 3D ?? ?? ?? ?? 04 75 4A D9 45 14 D8 0D ?? ?? ?? ?? D9 5D FC D9 45 FC E8 ?? ?? ?? ?? D9 45 10"),
			HL4554: hooks.MustMakePattern("51 83 3D ?? ?? ?? ?? 04 75 50 D9 44 24 14 D8 0D ?? ?? ?? ?? D9 5C 24 00 D9 44 24 00 E8 ?? ?? ?? ?? D9 44 24 10"),
		}),
		TriGLCullFace: hooks.MakeFunctionPattern("tri_GL_CullFace", nil, map[string]hooks.SearchPattern{
			HL8684: hooks.MustMakePattern("55 8B EC 8B 45 08 83 E8 00 74 10 48 75 23 68 44 0B 00 00 FF 15 ?? ?? ?? ?? 5D C3 68 44 0B 00 00"),
			HL4554: hooks.MustMakePattern("8B 44 24 04 83 E8 00 74 0F 48 75 22 68 44 0B 00 00 FF 15 ?? ?? ?? ?? C3 68 44 0B 00 00"),
		}),
		TriGLEnd: hooks.MakeFunctionPattern("tri_GL_End", nil, map[string]hooks.SearchPattern{
			HL8684: hooks.MustMakePattern("FF 25 ?? ?? ?? ?? 90 90 90 90 90 90 90 90 90 90 55 8B EC 8B 45 0C"),
			HL4554: hooks.MustMakePattern("FF 25 ?? ?? ?? ?? 90 90 90 90 90 90 90 90 90 90 8B 44 24 08 8B 4C 24 04"),
		}),
		TriGLRenderMode: hooks.MakeFunctionPattern("tri_GL_RenderMode", nil, map[string]hooks.SearchPattern{
			HL8684: hooks.MustMakePattern("55 8B EC 56 8B 75 08 83 FE 05 0F 87 ?? ?? ?? ?? FF 24 B5 ?? ?? ?? ?? 68 ?? ?? ?? ?? FF 15 ?? ?? ?? ?? 6A 01"),
			HL4554: hooks.MustMakePattern("56 8B 74 24 08 83 FE 05 0F 87 ?? ?? ?? ?? FF 24 B5 ?? ?? ?? ?? 68 ?? ?? ?? ?? FF 15 ?? ?? ?? ?? 6A 01"),
		}),
		TriGLVertex3fv: hooks.MakeFunctionPattern("tri_GL_Vertex3fv", nil, map[string]hooks.SearchPattern{
			HL8684: hooks.MustMakePattern("55 8B EC 8B 45 08 50 FF 15 ?? ?? ?? ?? 5D C3 90 55 8B EC 8B 45 10 8B 4D 0C 8B 55 08 50 51 52"),
			HL4554: hooks.MustMakePattern("8B 44 24 04 50 FF 15 ?? ?? ?? ?? C3 90 90 90 90 8B 44 24 0C 8B 4C 24 08 8B 54 24 04 50 51 52"),
		}),
		VFadeAlpha: hooks.MakeFunctionPattern("V_FadeAlpha", nil, map[string]hooks.SearchPattern{
			HL8684: hooks.MustMakePattern("55 8B EC 83 EC 08 D9 05 ?? ?? ?? ?? DC 1D ?? ?? ?? ?? 8A 0D ?? ?? ?? ?? DF E0 F6 C4 05 7A 1C D9 05 ?? ?? ?? ?? DC 1D"),
			HL4554: hooks.MustMakePattern("D9 05 ?? ?? ?? ?? DC 1D ?? ?? ?? ?? 8A 0D ?? ?? ?? ?? 83 EC 08 DF E0 F6 C4"),
		}),
		VGUI2DrawSetTextColorAlpha: hooks.MakeFunctionPattern("VGUI2_Draw_SetTextColorAlpha", nil, map[string]hooks.SearchPattern{
			HL8684: hooks.MustMakePattern("55 8B EC 8A 45 08 8A 4D 0C 8A 55 10 88 45 08 8A 45 14 88 4D 09 88 55 0A 88 45 0B 8B 4D 08 89"),
			HL4554: hooks.MustMakePattern("8A 44 24 04 8A 4C 24 08 8A 54 24 0C 88 44 24 04 8A 44 24 10 88 4C 24 05 88 54 24 06 88 44 24 07 8B 4C 24 04 89 0D"),
		}),
		WorldTransform: hooks.MakeFunctionPattern("WorldTransform", nil, map[string]hooks.SearchPattern{
			HL8684: hooks.MustMakePattern("55 8B EC 83 EC 08 8B 45 08 8B 4D 0C D9 05 ?? ?? ?? ?? D8 08 D9 05 ?? ?? ?? ?? D8 48"),
			HL4554: hooks.MustMakePattern("83 EC 08 8B 44 24 0C 8B 4C 24 10 D9 05 ?? ?? ?? ?? D8 08 D9 05 ?? ?? ?? ?? D8 48 08"),
		}),
		HUDDrawTransparentTriangles: hooks.MakeFunctionPattern("HUD_DrawTransparentTriangles", map[string]string{
			"Windows": "HUD_DrawTransparentTriangles",
		}, nil),
		HUDRedraw: hooks.MakeFunctionPattern("HUD_Redraw", map[string]string{
			"Windows": "HUD_Redraw",
		}, nil),
		HUDReset: hooks.MakeFunctionPattern("HUD_Reset", map[string]string{
			"Windows": "HUD_Reset",
		}, nil),
		HUDVidInit: hooks.MakeFunctionPattern("HUD_VidInit", map[string]string{
			"Windows": "HUD_VidInit",
		}, nil),
		CBaseMonsterChangeSchedule: hooks.MakeFunctionPattern("CBaseMonster::ChangeSchedule", map[string]string{
			WindowsHLDLL: "CBaseMonster::ChangeSchedule",
		}, map[string]hooks.SearchPattern{
			HL8684: hooks.MustMakePattern("8B 44 24 04 33 D2 89 81 78 01 00 00 89 91 7C 01 00 00 89 91 74 01 00 00 89 91 F0 00 00 00 89 91 68 02 00 00"),
			OF8684: hooks.MustMakePattern("8B 81 84 01 00 00 33 D2 3B C2 56 74 55 8B 00 3B C2 74 4F 8B B1 88 01 00 00 57 8B 3C F0"),
			HLWON:  hooks.MustMakePattern("8B 44 24 04 33 D2 89 81 44 01 00 00 89 91 48 01 00 00 89 91 40 01 00 00 89 91 BC 00 00 00 89 91 34 02 00 00"),
			OFWON:  hooks.MustMakePattern("8B 81 48 01 00 00 33 D2 3B C2 56 74 55 8B 00 3B C2 74 4F 8B B1 4C 01 00 00 57 8B 3C F0"),
			CSCZDS: hooks.MustMakePattern("8B 44 24 04 33 D2 89 81 74 01 00 00 89 91 78 01 00 00 89 91 70 01 00 00 89 91 7C 01 00 00 89 91 88 02 00 00"),
			Gunman: hooks.MustMakePattern("8B 44 24 04 53 57 8B F9 33 DB 89 87 4C 01 00 00 89 9F 50 01 00 00 89 9F 48 01 00 00 89 9F BC 00 00 00 89 9F 3C 02 00 00"),
		}),
		CBaseMonsterPBestSound: hooks.MakeFunctionPattern("CBaseMonster::PBestSound", map[string]string{
			WindowsHLDLL: "CBaseMonster::PBestSound",
		}, map[string]hooks.SearchPattern{
			HL8684: hooks.MustMakePattern("83 EC 10 53 8B D9 55 57 8B BB 1C 02 00 00 83 CD FF 83 FF FF C7 44 24 0C 00 00 00 46 75 2D"),
			OF8684: hooks.MustMakePattern("83 EC 10 53 8B D9 55 57 8B BB 28 02 00 00 83 CD FF 83 FF FF C7 44 24 0C 00 00 00 46 75 2D"),
			HLWON:  hooks.MustMakePattern("83 EC 10 53 8B D9 55 57 8B BB E8 01 00 00 83 CD FF 83 FF FF C7 44 24 0C 00 00 00 46 75 2D"),
			OFWON:  hooks.MustMakePattern("83 EC 10 53 8B D9 55 57 8B BB EC 01 00 00 83 CD FF 83 FF FF C7 44 24 0C 00 00 00 46 75 2D"),
			CSCZDS: hooks.MustMakePattern("83 EC 10 53 8B D9 55 57 8B BB 3C 02 00 00 83 CD FF 83 FF FF C7 44 24 0C 00 00 00 46 75 2D"),
			Gunman: hooks.MustMakePattern("83 EC 10 53 8B D9 55 57 8B BB F0 01 00 00 83 CD FF 83 FF FF C7 44 24 0C 00 00 00 46 75 2D"),
		}),
		CBaseMonsterRouteNew: hooks.MakeFunctionPattern("CBaseMonster::RouteNew", map[string]string{
			WindowsHLDLL: "CBaseMonster::RouteNew",
		}, map[string]hooks.SearchPattern{
			HL8684: hooks.MustMakePattern("33 C0 89 81 ?? ?? ?? ?? 89 81 ?? ?? ?? ?? C3 90 8B 81 ?? ?? ?? ?? C1 E0 04"),
			CSCZDS: hooks.MustMakePattern("33 C0 89 81 ?? ?? ?? ?? 89 81 ?? ?? ?? ?? C3 90 8B 81 ?? ?? ?? ?? 83 C0 14"),
			Gunman: hooks.MustMakePattern("33 C0 89 81 ?? ?? ?? ?? 89 81 ?? ?? ?? ?? C3 90 8B 81 ?? ?? ?? ?? 83 C0 16"),
		}),
		CGraphInitGraph: hooks.MakeFunctionPattern("CGraph::InitGraph", map[string]string{
			WindowsHLDLL: "CGraph::InitGraph",
		}, map[string]hooks.SearchPattern{
			HL8684: hooks.MustMakePattern("56 8B F1 57 33 FF 8B 46 10 89 3E 3B C7 89 7E 04 89 7E 08 74 0C 50 E8 ?? ?? ?? ?? 83 C4 04 89 7E 10 8B 46 0C"),
		}),
		CSoundEntActiveList: hooks.MakeFunctionPattern("CSoundEnt::ActiveList", map[string]string{
			WindowsHLDLL: "CSoundEnt::ActiveList",
		}, map[string]hooks.SearchPattern{
			HL8684: hooks.MustMakePattern("A1 ?? ?? ?? ?? 85 C0 75 04 83 C8 FF C3 8B 40 58 C3"),
			OF8684: hooks.MustMakePattern("A1 ?? ?? ?? ?? 85 C0 75 04 83 C8 FF C3 8B 40 64 C3"),
			HLWON:  hooks.MustMakePattern("A1 ?? ?? ?? ?? 85 C0 75 04 83 C8 FF C3 8B 40 24 C3"),
			CSCZDS: hooks.MustMakePattern("A1 ?? ?? ?? ?? 85 C0 75 04 83 C8 FF C3 8B 40 50 C3"),
		}),
		CSoundEntSoundPointerForIndex: hooks.MakeFunctionPattern("CSoundEnt::SoundPointerForIndex", map[string]string{
			WindowsHLDLL: "CSoundEnt::SoundPointerForIndex",
		}, map[string]hooks.SearchPattern{
			HL8684: hooks.MustMakePattern("8B 0D ?? ?? ?? ?? 85 C9 75 03 33 C0 C3 8B 44 24 04 83 F8 3F 7E 13 68 ?? ?? ?? ?? 6A 01 FF 15 ?? ?? ?? ??"),
		}),
		PMInit: hooks.MakeFunctionPattern("PM_Init", map[string]string{
			"Windows": "PM_Init",
		}, map[string]hooks.SearchPattern{
			HL8684: hooks.MustMakePattern("55 8B EC E8 ?? ?? ?? ?? 8B 55 08 33 C0 56 8D 8A ?? ?? ?? ?? 8B B0 ?? ?? ?? ?? 83 C0 0C 89 71 FC 8B B0 ?? ?? ?? ?? 89 31"),
			OF8684: hooks.MustMakePattern("8B 44 24 04 A3 ?? ?? ?? ?? E8 ?? ?? ?? ?? E8 ?? ?? ?? ?? C7 05"),
		}),
		PMPlayerMove: hooks.MakeFunctionPattern("PM_PlayerMove", map[string]string{
			"Windows": "PM_PlayerMove",
		}, map[string]hooks.SearchPattern{
			HL8684:     hooks.MustMakePattern("A1 ?? ?? ?? ?? 8B 4C 24 04 55 57 33 FF 89 48 04 E8 ?? ?? ?? ?? 8B 15 ?? ?? ?? ?? 33 C9 89 BA 8C 54 04 00 A1 ?? ?? ?? ?? 8A 88 5A 54 04 00 89"),
			BigLolly:   hooks.MustMakePattern("55 8B EC 83 EC 0C C7 45 FC 00 00 00 00 A1 ?? ?? ?? ?? 8B 4D 08 89 48 04 E8 ?? ?? ?? ?? 8B 15 ?? ?? ?? ?? C7 82 8C 54 04 00 00 00 00 00 A1"),
			TWHLTower2: hooks.MustMakePattern("55 8B EC 51 A1 ?? ?? ?? ?? 8B 4D 08 53 56 57 33 FF 89 7D FC 89 48 04 E8 D8 FC FF FF A1 ?? ?? ?? ?? 89 B8 8C 54 04 00 A1 ?? ?? ?? ?? 0F B6 88 5A 54 04 00"),
			CSCZDS:     hooks.MustMakePattern("A1 ?? ?? ?? ?? 8B 4C 24 04 55 56 57 33 ED 33 FF 89 48 04 E8 ?? ?? ?? ?? 8B 15 ?? ?? ?? ?? 33 C9 89 AA 8C 54 04 00 A1 ?? ?? ?? ?? 8A 88 5A 54 04 00 89"),
		}),
		WorldGraph: hooks.MakeFunctionPattern("WorldGraph", map[string]string{
			WindowsHLDLL: "WorldGraph",
		}, nil),
	}
}
