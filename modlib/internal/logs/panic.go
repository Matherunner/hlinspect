package logs

import "runtime/debug"

func HandlePanic() {
	if r := recover(); r != nil {
		DLLLog.Errorf("got a panic\nr: %+v\nstack: %+v", r, string(debug.Stack()))
		panic(r)
	}
}
