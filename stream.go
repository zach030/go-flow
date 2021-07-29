package go_flow

type InPoint interface {
	In() chan<-interface{}
}

type OutPoint interface {
	Out() <-chan interface{}
}

type Source interface {
	OutPoint
}

type Target interface {
	InPoint
}



