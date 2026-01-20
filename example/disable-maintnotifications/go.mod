module github.com/hanzoai/kv-go/example/disable-maintnotifications

go 1.23

replace github.com/hanzoai/kv-go/v9 => ../..

require github.com/hanzoai/kv-go/v9 v9.7.0

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	go.uber.org/atomic v1.11.0 // indirect
)
