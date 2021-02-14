package hw

import (
	"hlinspect/internal/hooks"
	"hlinspect/internal/logs"
)

var hwDLL *hooks.Module

var buildNumberPattern = hooks.MakeFunctionPattern("build_number", map[string]hooks.SearchPattern{
	"HL-SteamPipe": hooks.MustMakePattern("55 8B EC 83 EC 08 A1 ?? ?? ?? ?? 56 33 F6 85 C0 0F 85 9B 00 00 00 53 33 DB 8B 04 9D ?? ?? ?? ?? 8B 0D ?? ?? ?? ?? 6A 03 50 51 E8"),
})

// BuildNumber build_number
func BuildNumber() int {
	return hooks.CallFuncInts0(buildNumberPattern.GetAddress())
}

// InitHWDLL initialise hw.dll hooks and symbol search
func InitHWDLL() (err error) {
	if hwDLL != nil {
		return
	}

	hwDLL, err = hooks.NewModule("hw.dll")
	if err != nil {
		return
	}

	name, addr, err := buildNumberPattern.Find(hwDLL)
	logs.DLLLog.Debugf("Build number search: %v %v %v %v\n", name, addr, err, BuildNumber())

	return
}
