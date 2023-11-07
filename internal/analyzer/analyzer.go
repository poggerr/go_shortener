package analyzer

import (
	"golang.org/x/tools/go/analysis"
)

var OsExitCheckAnalyzer = &analysis.Analyzer{
	Name: "check os.Exit",
	Doc:  "check for os.Exit calls",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	// дописать свой анализатор
	return nil, nil
}
