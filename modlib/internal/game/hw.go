package game

import (
	"hlinspect/internal/engine"
	"hlinspect/internal/game/cdefs"
	"hlinspect/internal/game/registry"
	"hlinspect/internal/hooks"
	"hlinspect/internal/logs"
	"unsafe"
)

var hwDLL *hooks.Module

func initHWDLL(base string) (err error) {
	if hwDLL != nil {
		return
	}

	hwDLL, err = hooks.NewModule(base)
	if err != nil {
		return
	}

	reg := Model.Registry()

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
		&reg.CbufInsertText:             nil,
	}

	errors := hooks.BatchFind(hwDLL, items)
	printBatchFindErrors(errors)

	switch reg.HostAutoSaveF.PatternKey() {
	case registry.VersionHL8684:
		ptr := *(*unsafe.Pointer)(unsafe.Add(reg.HostAutoSaveF.Ptr(), 19))
		engine.Engine.SetSV(ptr)
		logs.DLLLog.Debugf("Set SV address: %x", ptr)
	}

	switch reg.HostNoclipF.PatternKey() {
	case registry.VersionHL8684:
		ptr := *(*unsafe.Pointer)(unsafe.Add(reg.HostNoclipF.Ptr(), 31))
		ptr = unsafe.Add(ptr, -0x14)
		engine.Engine.SetGlobalVariables(ptr)
		logs.DLLLog.Debugf("Set GlobalVariables address: %x", ptr)
	case registry.VersionHL4554:
		ptr := *(*unsafe.Pointer)(unsafe.Add(reg.HostNoclipF.Ptr(), 28))
		ptr = unsafe.Add(ptr, -0x14)
		engine.Engine.SetGlobalVariables(ptr)
		logs.DLLLog.Debugf("Set GlobalVariables address: %x", ptr)
	case registry.VersionHLNGHL:
		ptr := *(*unsafe.Pointer)(unsafe.Add(reg.HostNoclipF.Ptr(), 27))
		ptr = unsafe.Add(ptr, -0x14)
		engine.Engine.SetGlobalVariables(ptr)
		logs.DLLLog.Debugf("Set GlobalVariables address: %x", ptr)
	}

	return nil
}
