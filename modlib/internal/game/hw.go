package game

import (
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
		&reg.WriteDestParm:              nil,
		&reg.SVExecuteClientMessage:     cdefs.CDefs.HookedSVExecuteClientMessage,
	}

	errors := hooks.BatchFind(hwDLL, items)
	printBatchFindErrors(errors)

	initGlobalSV()

	switch reg.HostNoclipF.PatternKey() {
	case registry.VersionHL8684:
		ptr := *(*unsafe.Pointer)(unsafe.Add(reg.HostNoclipF.Ptr(), 31))
		ptr = unsafe.Add(ptr, -0x14)
		Model.S().SetGlobalVariables(ptr)
		logs.DLLLog.Debugf("Set GlobalVariables address: %x", ptr)
	case registry.VersionHL4554:
		ptr := *(*unsafe.Pointer)(unsafe.Add(reg.HostNoclipF.Ptr(), 28))
		ptr = unsafe.Add(ptr, -0x14)
		Model.S().SetGlobalVariables(ptr)
		logs.DLLLog.Debugf("Set GlobalVariables address: %x", ptr)
	case registry.VersionHLNGHL:
		ptr := *(*unsafe.Pointer)(unsafe.Add(reg.HostNoclipF.Ptr(), 27))
		ptr = unsafe.Add(ptr, -0x14)
		Model.S().SetGlobalVariables(ptr)
		logs.DLLLog.Debugf("Set GlobalVariables address: %x", ptr)
	}

	return nil
}

func initGlobalSV() {
	// WriteDest_Parm returns the address of sv.datagram if the argument is 0.
	ptr := Model.API().WriteDestParm(0)
	// Offset to get the address of sv
	ptr = unsafe.Add(ptr, -0x3bc68)
	Model.S().SetSV(ptr)
	logs.DLLLog.Debugf("Set SV address: %x", ptr)
}
