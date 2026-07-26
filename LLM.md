# LLM.md — Hanzo KV, Go client

Module: `github.com/hanzokv/go/v9` · package `kv`

```bash
go build ./...
go test ./...      # the integration suites need a live KV on :6379
```

## The vocabulary is `kv`

We do not use Redis. The client says so end to end:

| | |
|---|---|
| package | `kv` — `kv.NewClient`, `kv.Options`, `kv.Nil`, `kv.UniversalClient` |
| errors | `"kv: …"` |
| URL scheme | **`kv://`, `kvs://`, `unix://`** — `redis://` does not parse |
| companions | `extra/kv{otel,cmd,prometheus,census}` |
| env | `KV_VERSION`, `KV_ENDPOINTS_CONFIG_PATH` |

`ParseURL` accepting only `kv://` fails closed by design. A stale
`KV_URL=redis://…` stops a service that validates it and silently degrades one
that does not. No alias, no dual-accept.

Two spellings stay, and neither is branding: **`RediSearch`** names a specific
external module the `FT.*` commands map to, and **`StoreDist`** matches a
case-insensitive search for "redis" by accident — `sto-reDis-t`.

None of this touches the wire. RESP commands are `GET`/`SET`/`HSET`; the word
never appears in the protocol. The only server-visible spellings were the
`ParseURL` scheme, which we choose, and the `CLIENT SETINFO` lib-name label,
which the server merely records.

## One of each thing in the org

| repo | is |
|---|---|
| `hanzokv/go` | this client |
| `hanzokv/js` | the TypeScript client — npm **`@hanzo/kv`** |
| `hanzokv/mock` | the test mock — typed against this client, so it moves in lockstep |
| `hanzokv/lock` | distributed locks — its adapter is typed against this client too |
| `hanzokv/valkey` | the server fork |

Named for the language or the concern; the org supplies the product noun. A repo
called `kv-go` under `hanzokv` would say it twice. Same convention as `hanzo-ds`.

## Versioning

Always `x.y.z` → `x.y.z+1`. This module stays **v9** forever — breaking changes
are minor bumps, so there is no `/v10` import churn.

**v9.21.1 is the floor for the module path** and **v9.22.0 for the package
rename**. Anything below declares an older path or the old package name, and Go
rejects a module whose declared path does not match the requested one.

`extra/*` are separate modules with their own tags (`extra/kvotel/v9.22.0`) and
carry `replace` directives to the parent for local work, so a release is the core
tag plus four extras.
