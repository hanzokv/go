# Disable Maintenance Notifications Example

This example demonstrates how to use the go-kv client with maintenance notifications **disabled**.

## What are Maintenance Notifications?

Maintenance notifications are a KV Cloud feature that allows the server to notify clients about:
- Planned maintenance events
- Failover operations
- Node migrations
- Cluster topology changes

The go-kv client supports three modes:
- **`ModeDisabled`**: Client doesn't send `CLIENT MAINT_NOTIFICATIONS ON` command
- **`ModeEnabled`**: Client forcefully sends the command, interrupts connection on error
- **`ModeAuto`** (default): Client tries to send the command, disables feature on error

## When to Disable Maintenance Notifications

You should disable maintenance notifications when:

1. **Connecting to non-KV Cloud / KV Enterprise instances** - Standard KV servers don't support this feature
2. **You want to handle failovers manually** - Your application has custom failover logic
3. **Minimizing client-side overhead** - You want the simplest possible client behavior
4. **The KV server doesn't support the feature** - Older KV versions or forks

## Usage

### Basic Example

```go
import (
    "github.com/kv/go-kv/v9"
    "github.com/kv/go-kv/v9/maintnotifications"
)

rdb := kv.NewClient(&kv.Options{
    Addr: "localhost:6379",

    // Explicitly disable maintenance notifications
    MaintNotificationsConfig: &maintnotifications.Config{
        Mode: maintnotifications.ModeDisabled,
    },
})
defer rdb.Close()
```

### Cluster Client Example

```go
rdbCluster := kv.NewClusterClient(&kv.ClusterOptions{
    Addrs: []string{"localhost:7000", "localhost:7001", "localhost:7002"},

    // Disable maintenance notifications for cluster
    MaintNotificationsConfig: &maintnotifications.Config{
        Mode: maintnotifications.ModeDisabled,
    },
})
defer rdbCluster.Close()
```

### Default Behavior (ModeAuto)

If you don't specify `MaintNotifications`, the client defaults to `ModeAuto`:

```go
// This uses ModeAuto by default
rdb := kv.NewClient(&kv.Options{
    Addr: "localhost:6379",
    // MaintNotificationsConfig: nil means ModeAuto
})
```

With `ModeAuto`, the client will:
1. Try to enable maintenance notifications
2. If the server doesn't support it, silently disable the feature
3. Continue normal operation

## Running the Example

1. Start a KV server:
   ```bash
   kv-server --port 6379
   ```

2. Run the example:
   ```bash
   go run main.go
   ```

## Expected Output

```
=== Example 1: Explicitly Disabled ===
✓ Connected successfully (maintenance notifications disabled)
✓ SET operation successful
✓ GET operation successful: value1

=== Example 2: Default Behavior (ModeAuto) ===
✓ Connected successfully (maintenance notifications auto-enabled)

=== Example 3: Cluster Client with Disabled Notifications ===
Cluster not available (expected): ...

=== Example 4: Performance Comparison ===
✓ 1000 SET operations (disabled): 45ms
✓ 1000 SET operations (auto): 46ms

=== Cleanup ===
✓ Database flushed

=== Summary ===
Maintenance notifications can be disabled by setting:
  MaintNotificationsConfig: &maintnotifications.Config{
    Mode: maintnotifications.ModeDisabled,
  }

This is useful when:
  - Connecting to non-KV Cloud instances
  - You want to handle failovers manually
  - You want to minimize client-side overhead
  - The KV server doesn't support CLIENT MAINT_NOTIFICATIONS
```

## Performance Impact

Disabling maintenance notifications has minimal performance impact. The main differences are:

1. **Connection Setup**: One less command (`CLIENT MAINT_NOTIFICATIONS ON`) during connection initialization
2. **Runtime Overhead**: No background processing of maintenance notifications
3. **Memory Usage**: Slightly lower memory footprint (no notification handlers)

In most cases, the performance difference is negligible (< 1%).