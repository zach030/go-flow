package flow

import goflow "go-flow"

type MapFunc func(interface{}) interface{}

type Map struct {
	fun         MapFunc
	in          chan interface{}
	out         chan interface{}
	parallelism uint
}

func NewMap(fun MapFunc, parallelism uint) *Map {
	m := &Map{
		fun:         fun,
		in:          make(chan interface{}),
		out:         make(chan interface{}),
		parallelism: parallelism,
	}
	go m.doMap()
	return m
}

func (m *Map) doMap() {
	sem := make(chan struct{}, m.parallelism)
	for elem := range m.in {
		sem <- struct{}{}
		go func(e interface{}) {
			defer func() { <-sem }()
			transData := m.fun(e)
			m.out <- transData
		}(elem)
	}
	for i := 0; i < int(m.parallelism); i++ {
		sem <- struct{}{}
	}
	close(m.out)
}

func (m *Map) Input() chan<- interface{} {
	return m.in
}

func (m *Map) Output() <-chan interface{} {
	return m.out
}

func (m *Map) Via(flow goflow.DataFlow) goflow.DataFlow {
	go m.transmit(flow)
	return flow
}

func (m *Map) transmit(in goflow.InPoint) {
	for elem := range m.Output() {
		in.Input() <- elem
	}
	close(in.Input())
}

func (m *Map) To(target goflow.Target) {
	m.transmit(target)
}
