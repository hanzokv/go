module github.com/hanzokv/go/extra/rediscensus/v9

go 1.26

replace github.com/hanzokv/go/v9 => ../..

replace github.com/hanzokv/go/extra/rediscmd/v9 => ../rediscmd

require (
	github.com/hanzokv/go/extra/rediscmd/v9 v9.21.1
	github.com/hanzokv/go/v9 v9.21.1
	go.opencensus.io v0.24.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	go.uber.org/atomic v1.11.0 // indirect
)

retract (
	v9.7.2 // This version was accidentally released. Please use version 9.7.3 instead.
	v9.5.3 // This version was accidentally released. Please use version 9.6.0 instead.
)
