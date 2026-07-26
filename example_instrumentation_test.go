package kv_test

import (
	"context"
	"fmt"
	"net"

	"github.com/hanzokv/go/v9"
)

type kvHook struct{}

var _ kv.Hook = kvHook{}

func (kvHook) DialHook(hook kv.DialHook) kv.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		fmt.Printf("dialing %s %s\n", network, addr)
		conn, err := hook(ctx, network, addr)
		fmt.Printf("finished dialing %s %s\n", network, addr)
		return conn, err
	}
}

func (kvHook) ProcessHook(hook kv.ProcessHook) kv.ProcessHook {
	return func(ctx context.Context, cmd kv.Cmder) error {
		fmt.Printf("starting processing: <%v>\n", cmd.Args())
		err := hook(ctx, cmd)
		fmt.Printf("finished processing: <%v>\n", cmd.Args())
		return err
	}
}

func (kvHook) ProcessPipelineHook(hook kv.ProcessPipelineHook) kv.ProcessPipelineHook {
	return func(ctx context.Context, cmds []kv.Cmder) error {
		names := make([]string, 0, len(cmds))
		for _, cmd := range cmds {
			names = append(names, fmt.Sprintf("%v", cmd.Args()))
		}
		fmt.Printf("pipeline starting processing: %v\n", names)
		err := hook(ctx, cmds)
		fmt.Printf("pipeline finished processing: %v\n", names)
		return err
	}
}

func Example_instrumentation() {
	rdb := kv.NewClient(&kv.Options{
		Addr:            ":6379",
		DisableIdentity: true,
	})
	rdb.AddHook(kvHook{})

	rdb.Ping(ctx)
	// Output:
	// starting processing: <[ping]>
	// dialing tcp :6379
	// finished dialing tcp :6379
	// starting processing: <[hello 3]>
	// finished processing: <[hello 3]>
	// starting processing: <[client maint_notifications on moving-endpoint-type internal-fqdn]>
	// finished processing: <[client maint_notifications on moving-endpoint-type internal-fqdn]>
	// finished processing: <[ping]>
}

func ExamplePipeline_instrumentation() {
	rdb := kv.NewClient(&kv.Options{
		Addr:            ":6379",
		DisableIdentity: true,
	})
	rdb.AddHook(kvHook{})

	rdb.Pipelined(ctx, func(pipe kv.Pipeliner) error {
		pipe.Ping(ctx)
		pipe.Ping(ctx)
		return nil
	})
	// Output:
	// pipeline starting processing: [[ping] [ping]]
	// dialing tcp :6379
	// finished dialing tcp :6379
	// starting processing: <[hello 3]>
	// finished processing: <[hello 3]>
	// starting processing: <[client maint_notifications on moving-endpoint-type internal-fqdn]>
	// finished processing: <[client maint_notifications on moving-endpoint-type internal-fqdn]>
	// pipeline finished processing: [[ping] [ping]]
}

func ExampleClient_Watch_instrumentation() {
	rdb := kv.NewClient(&kv.Options{
		Addr:            ":6379",
		DisableIdentity: true,
	})
	rdb.AddHook(kvHook{})

	rdb.Watch(ctx, func(tx *kv.Tx) error {
		tx.Ping(ctx)
		tx.Ping(ctx)
		return nil
	}, "foo")
	// Output:
	// starting processing: <[watch foo]>
	// dialing tcp :6379
	// finished dialing tcp :6379
	// starting processing: <[hello 3]>
	// finished processing: <[hello 3]>
	// starting processing: <[client maint_notifications on moving-endpoint-type internal-fqdn]>
	// finished processing: <[client maint_notifications on moving-endpoint-type internal-fqdn]>
	// finished processing: <[watch foo]>
	// starting processing: <[ping]>
	// finished processing: <[ping]>
	// starting processing: <[ping]>
	// finished processing: <[ping]>
	// starting processing: <[unwatch]>
	// finished processing: <[unwatch]>
}
