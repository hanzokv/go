module github.com/hanzoai/kv-go/example/redis-bloom

go 1.21

replace github.com/hanzoai/kv-go/v9 => ../..

require github.com/hanzoai/kv-go/v9 v9.18.0-beta.2

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	go.uber.org/atomic v1.11.0 // indirect
)
