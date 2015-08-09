package trace

import (
	"io"
	"fmt"
)

//	Tracer is the interface that describes an object capable of tracing events throughout code.
//	Capital letters mean this is a publicly visible type
type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a...interface{}) {
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}

func New(w io.Writer) Tracer {
	return nil
}