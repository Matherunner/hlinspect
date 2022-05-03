package gamelibs

import (
	"fmt"
	"hlinspect/internal/hooks"
	"hlinspect/internal/logs"
)

func printBatchFindErrors(errors map[*hooks.FunctionPattern]error) {
	for pat, err := range errors {
		if err == nil {
			useType := ""
			if pat.PatternKey() != "" {
				useType = fmt.Sprintf("pattern %v", pat.PatternKey())
			} else if pat.SymbolKey() != "" {
				useType = fmt.Sprintf("symbol %v", pat.SymbolKey())
			}
			logs.DLLLog.Debugf("Found %v at %v using %v", pat.Name(), pat.Address(), useType)
		} else {
			logs.DLLLog.Debugf("Failed to find %v: %v", pat.Name(), err)
		}
	}
}
