module github.com/hanzoai/kv-go/example/cluster-mget

go 1.26

replace github.com/hanzoai/kv-go/v9 => ../..

require github.com/hanzoai/kv-go/v9 v9.16.0

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
)
