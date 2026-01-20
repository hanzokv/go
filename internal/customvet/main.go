package main

import (
	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/hanzoai/kv-go/internal/customvet/checks/setval"
)

func main() {
	multichecker.Main(
		setval.Analyzer,
	)
}
