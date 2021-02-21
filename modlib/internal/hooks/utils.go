package hooks

import (
	"unsafe"
)

// BatchFind find and optionally hook a list of items in parallel
func BatchFind(mod *Module, items map[*FunctionPattern]unsafe.Pointer) (errors map[*FunctionPattern]error) {
	worker := func(pattern *FunctionPattern, hookFunc unsafe.Pointer, ch chan<- error) {
		if hookFunc == nil {
			_, _, err := pattern.Find(mod)
			ch <- err
		} else {
			_, _, err := pattern.Hook(mod, hookFunc)
			ch <- err
		}
		close(ch)
	}

	chans := make(map[*FunctionPattern]chan error)
	for pat, hookFunc := range items {
		chans[pat] = make(chan error)
		go worker(pat, hookFunc, chans[pat])
	}

	errors = make(map[*FunctionPattern]error)
	for pat, ch := range chans {
		errors[pat] = <-ch
	}

	return
}
