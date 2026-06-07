module github.com/hanzoai/kv-go/extra/rediscmd/v9

go 1.26.4

replace github.com/hanzoai/kv-go/v9 => ../..

require (
	github.com/bsm/ginkgo/v2 v2.12.0
	github.com/bsm/gomega v1.27.10
	github.com/hanzoai/kv-go/v9 v9.18.0-beta.2
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
)

retract (
	v9.7.2 // This version was accidentally released. Please use version 9.7.3 instead.
	v9.5.3 // This version was accidentally released. Please use version 9.6.0 instead.
)
