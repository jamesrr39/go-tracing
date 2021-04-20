package tracing

import (
	"context"
	fmt "fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	_ "github.com/gogo/protobuf/gogoproto" // for protobuf generation libraries to be included
	proto "github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
)

type key int

var (
	TraceCtxKey  key = 1
	TracerCtxKey key = 2
)

func Middleware(tracer *Tracer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			traceName := fmt.Sprintf("%s: %s", r.URL.String(), uuid.New().String())

			trace := StartTrace(tracer, traceName)

			newCtx := r.Context()
			newCtx = context.WithValue(newCtx, TraceCtxKey, trace)
			newCtx = context.WithValue(newCtx, TracerCtxKey, tracer)

			r = r.WithContext(newCtx)

			next.ServeHTTP(w, r)

			err := tracer.EndTrace(trace, "")
			if err != nil {
				log.Printf("Tracer: could not EndTrace. Error: %q\n", err)
			}
		}

		return http.HandlerFunc(fn)
	}
}

type Tracer struct {
	writer  io.Writer
	nowFunc func() time.Time
	writeMu sync.Mutex
}

func NewTracer(writer io.Writer) *Tracer {
	return &Tracer{writer, time.Now, sync.Mutex{}}
}

func StartTrace(tracer *Tracer, name string) *Trace {
	return &Trace{
		Name:           name,
		StartTimeNanos: tracer.nowFunc().UnixNano(),
		Spans:          []*Span{},
	}
}

func StartSpan(ctx context.Context, name string) (*Span, error) {
	tracerVal := ctx.Value(TracerCtxKey)
	if tracerVal == nil {
		return nil, fmt.Errorf("Trace: no tracer in context")
	}

	tracer := tracerVal.(*Tracer)

	return &Span{
		Name:           name,
		StartTimeNanos: tracer.nowFunc().UnixNano(),
	}, nil
}

func (span *Span) End(ctx context.Context) error {
	traceVal := ctx.Value(TraceCtxKey)
	if traceVal == nil {
		return fmt.Errorf("Trace: no trace in context")
	}

	tracerVal := ctx.Value(TracerCtxKey)
	if tracerVal == nil {
		return fmt.Errorf("Trace: no tracer in context")
	}

	trace := traceVal.(*Trace)
	tracer := tracerVal.(*Tracer)
	span.EndTimeNanos = tracer.nowFunc().UnixNano()

	trace.Spans = append(trace.Spans, span)
	return nil
}

func (tracer *Tracer) EndTrace(trace *Trace, summary string) error {

	endTime := tracer.nowFunc()

	trace.EndTimeNanos = endTime.UnixNano()
	trace.Summary = summary

	b, err := proto.Marshal(trace)
	if err != nil {
		return fmt.Errorf("Trace: error marshalling trace to protobuf message: %q", err)
	}

	tracer.writeMu.Lock()
	defer tracer.writeMu.Unlock()

	_, err = tracer.writer.Write(b)
	if err != nil {
		return fmt.Errorf("Trace: error writing trace to file/writable: %q", err)
	}

	return nil
}
