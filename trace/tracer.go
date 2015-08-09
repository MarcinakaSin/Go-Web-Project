package trace
//	Tracer is the interface that describes an object capable of tracing events throughout code.
//	Capital letters mean this is a publicly visible type
type Tracer interface {
	Trace(...interface{})
}