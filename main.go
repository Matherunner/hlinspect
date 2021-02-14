package main

/*
#include "dllmain.h"
*/
import "C"

func main() {}

// OnProcessAttach called from DllMain on process attach
//export OnProcessAttach
func OnProcessAttach() {
	initHooks()
}

// OnProcessDetach called from DllMain on process detach
//export OnProcessDetach
func OnProcessDetach() {
	cleanupHooks()
}
