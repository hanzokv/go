module github.com/hanzoai/kv-go/example/digest-optimistic-locking

go 1.26.3

replace github.com/hanzoai/kv-go/v9 => ../..

require github.com/hanzoai/kv-go/v9 v9.18.0-beta.2

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	github.com/zeebo/xxh3 v1.1.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
)
