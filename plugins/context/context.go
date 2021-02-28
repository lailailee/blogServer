package context

type Context interface {
	// get current goroutine cid.
	Cid() int
}

// implements the context.
// @remark user can use nil context.
type context int

var __cid int = 100

func NewContext() Context {
	v := context(__cid)
	__cid++
	return v
}

func (v context) Cid() int {
	return int(v)
}
