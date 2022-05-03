package registry

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"hlinspect/internal/hooks"
	"reflect"
)

var (
	ErrFieldNoTag         = errors.New("struct field has no tag")
	ErrInvalidPatternName = errors.New("invalid pattern name for struct field")
)

//go:embed patterns/patterns.json
var patternsJSON []byte

type patternsRoot struct {
	Symbols map[string]symbolDefinition `json:"symbols"`
}

type symbolDefinition struct {
	Comment  string            `json:"#"`
	DLLs     []string          `json:"dlls"`
	Names    map[string]string `json:"names"`
	Patterns map[string]string `json:"patterns"`
}

// API holds the addresses to game DLL functions.
type API struct {
	// HW
	AngleVectors               hooks.FunctionPattern `patname:"AngleVectors"`
	BuildNumber                hooks.FunctionPattern `patname:"build_number"`
	CmdAddCommandWithFlags     hooks.FunctionPattern `patname:"Cmd_AddCommandWithFlags"`
	CmdArgv                    hooks.FunctionPattern `patname:"Cmd_Argv"`
	CvarRegisterVariable       hooks.FunctionPattern `patname:"Cvar_RegisterVariable"`
	DrawString                 hooks.FunctionPattern `patname:"Draw_String"`
	HostAutoSaveF              hooks.FunctionPattern `patname:"Host_AutoSave_f"`
	HostNoclipF                hooks.FunctionPattern `patname:"Host_Noclip_f"`
	HudGetScreenInfo           hooks.FunctionPattern `patname:"hudGetScreenInfo"`
	MemoryInit                 hooks.FunctionPattern `patname:"Memory_Init"`
	PFCheckClientI             hooks.FunctionPattern `patname:"PF_checkclient_I"`
	PFTracelineDLL             hooks.FunctionPattern `patname:"PF_traceline_DLL"`
	RClear                     hooks.FunctionPattern `patname:"R_Clear"`
	RDrawSequentialPoly        hooks.FunctionPattern `patname:"R_DrawSequentialPoly"`
	ScreenTransform            hooks.FunctionPattern `patname:"ScreenTransform"`
	TriGLBegin                 hooks.FunctionPattern `patname:"tri_GL_Begin"`
	TriGLColor4f               hooks.FunctionPattern `patname:"tri_GL_Color4f"`
	TriGLCullFace              hooks.FunctionPattern `patname:"tri_GL_CullFace"`
	TriGLEnd                   hooks.FunctionPattern `patname:"tri_GL_End"`
	TriGLRenderMode            hooks.FunctionPattern `patname:"tri_GL_RenderMode"`
	TriGLVertex3fv             hooks.FunctionPattern `patname:"tri_GL_Vertex3fv"`
	VFadeAlpha                 hooks.FunctionPattern `patname:"V_FadeAlpha"`
	VGUI2DrawSetTextColorAlpha hooks.FunctionPattern `patname:"VGUI2_Draw_SetTextColorAlpha"`
	WorldTransform             hooks.FunctionPattern `patname:"WorldTransform"`

	// CL
	HUDDrawTransparentTriangles hooks.FunctionPattern `patname:"HUD_DrawTransparentTriangles"`
	HUDRedraw                   hooks.FunctionPattern `patname:"HUD_Redraw"`
	HUDReset                    hooks.FunctionPattern `patname:"HUD_Reset"`
	HUDVidInit                  hooks.FunctionPattern `patname:"HUD_VidInit"`

	// HL
	CBaseMonsterChangeSchedule    hooks.FunctionPattern `patname:"CBaseMonster::ChangeSchedule"`
	CBaseMonsterPBestSound        hooks.FunctionPattern `patname:"CBaseMonster::PBestSound"`
	CBaseMonsterRouteNew          hooks.FunctionPattern `patname:"CBaseMonster::RouteNew"`
	CGraphInitGraph               hooks.FunctionPattern `patname:"CGraph::InitGraph"`
	CSoundEntActiveList           hooks.FunctionPattern `patname:"CSoundEnt::ActiveList"`
	CSoundEntSoundPointerForIndex hooks.FunctionPattern `patname:"CSoundEnt::SoundPointerForIndex"`
	PMInit                        hooks.FunctionPattern `patname:"PM_Init"`
	PMPlayerMove                  hooks.FunctionPattern `patname:"PM_PlayerMove"`
	WorldGraph                    hooks.FunctionPattern `patname:"WorldGraph"`
}

func NewAPI() (*API, error) {
	var patterns patternsRoot
	err := json.Unmarshal(patternsJSON, &patterns)
	if err != nil {
		return nil, err
	}

	api := &API{}

	v := reflect.Indirect(reflect.ValueOf(api))
	for i := 0; i < v.NumField(); i++ {
		patName, ok := v.Type().Field(i).Tag.Lookup("patname")
		if !ok {
			return nil, ErrFieldNoTag
		}

		pat, ok := patterns.Symbols[patName]
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrInvalidPatternName, patName)
		}

		patterns := map[string]hooks.SearchPattern{}
		for k, v := range pat.Patterns {
			patterns[k] = hooks.MustPattern(v)
		}

		patVal := hooks.NewFunctionPattern(patName, pat.Names, patterns)
		reflPatVal := reflect.ValueOf(patVal)

		fieldVal := v.Field(i)
		fieldVal.Set(reflPatVal)
	}

	return api, nil
}
