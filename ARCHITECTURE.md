# Architecture

## Overview

Like many hooking-based mods, the features in HLInspect are built on top of the underlying game DLL functions. For example, registering a custom cvar requires calling into the game DLL. Similarly, to run some logic when the game decides to refresh the HUD, HLInspect has to hook into the `HUD_Redraw` function and run the custom HUD features when the hooked function is called. In general, HLInspect has no real control over when the hooked functions will be called. This is nothing more than an **event-driven architecture**. In light of this, we can view all hooked functions to be merely callbacks to events coming from the game DLLs. This allows us to borrow experience of event-driven programming from other domains such as GUI programming, along with all the solutions, caveats, and pitfalls of this development methodology.

The bottom layer of HLInspect is the code doing the dirty and messy work of interfacing with the game DLLs and receiving events via the hooked functions. This part of the code is confined to the `gamelibs` folder. This layer exposes a slightly higher-level and more idiomatic API to other parts of HLInspect. This layer can be viewed as a black box with functions to call into and receive events from. As a general rule, the exposed APIs should never use C types. But `unsafe.Pointer` and `uintptr` are permitted when appropriate.

The main entry point of HLInspect initialises and hooks the bottom layer, and registers a default event listener. The rest of the code is largely built on top of the event listener and the API provided by `gamelibs`.

## Hooking

There are three ways to find a symbol. The first is to rely on the symbol name. This only works if the symbol is explicitly exported from the DLL, such as `HUD_Redraw`. On Windows, this is done by calling the `GetProcAddress` API on each game DLL module. If this fails, HLInspect will try to search for the symbol using debug information. This is done using the DbgHelp family of Windows API. This is extremely useful for a researcher studying Half-Life using custom DLLs built from Half-Life SDK in Debug mode, while retaining the ability to use the Visual Studio debugger on the Half-Life process. When all else fails, HLInspect will resort to signature or pattern searching, which is a standard technique for mods of this nature. The symbol finding code is parallelised using goroutines.

One the symbols have been found, HLInspect will proceed hook into some of them, as defined by the programmer. HLInspect relies on the famous [MinHook](https://github.com/TsudaKageyu/minhook) library for its hooking capability.
