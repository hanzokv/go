package kvcensus

import (
	"context"
	"net"

	"go.opencensus.io/trace"

	"github.com/hanzokv/go/extra/kvcmd/v9"
	"github.com/hanzokv/go/v9"
)

type TracingHook struct{}

var _ kv.Hook = (*TracingHook)(nil)

func NewTracingHook() *TracingHook {
	return new(TracingHook)
}

func (TracingHook) DialHook(next kv.DialHook) kv.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		ctx, span := trace.StartSpan(ctx, "dial")
		defer span.End()

		span.AddAttributes(
			trace.StringAttribute("db.system", "kv"),
			trace.StringAttribute("network", network),
			trace.StringAttribute("addr", addr),
		)

		conn, err := next(ctx, network, addr)
		if err != nil {
			recordErrorOnOCSpan(ctx, span, err)

			return nil, err
		}

		return conn, nil
	}
}

func (TracingHook) ProcessHook(next kv.ProcessHook) kv.ProcessHook {
	return func(ctx context.Context, cmd kv.Cmder) error {
		ctx, span := trace.StartSpan(ctx, cmd.FullName())
		defer span.End()

		span.AddAttributes(
			trace.StringAttribute("db.system", "kv"),
			trace.StringAttribute("kv.cmd", kvcmd.CmdString(cmd)),
		)

		err := next(ctx, cmd)
		if err != nil {
			recordErrorOnOCSpan(ctx, span, err)
			return err
		}

		if err = cmd.Err(); err != nil {
			recordErrorOnOCSpan(ctx, span, err)
		}

		return nil
	}
}

func (TracingHook) ProcessPipelineHook(next kv.ProcessPipelineHook) kv.ProcessPipelineHook {
	return next
}

func recordErrorOnOCSpan(ctx context.Context, span *trace.Span, err error) {
	if err != kv.Nil {
		span.AddAttributes(trace.BoolAttribute("error", true))
		span.Annotate([]trace.Attribute{trace.StringAttribute("Error", "kv error")}, err.Error())
	}
}
