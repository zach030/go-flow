package go_flow

// 接收数据方法
type InPoint interface {
	Input() chan<- interface{}
}

// 输出数据方
type OutPoint interface {
	Output() <-chan interface{}
}

// 数据源
type Source interface {
	OutPoint
	Via(DataFlow) DataFlow
}

// 数据目的地
type Target interface {
	InPoint
	OutPoint
	//Call(hook Hook)
}

// 传输数据流
type DataFlow interface {
	InPoint
	OutPoint
	Via(DataFlow) DataFlow
	To(Target)
	//Call(hook Hook)
}

// 数据处理结束之后
type Hook interface {
	InPoint
	Resp() interface{}
	Result() interface{}
}

func FlowTo(out OutPoint, in InPoint) {
	go func() {
		for elem := range out.Output() {
			in.Input() <- elem
		}
		close(in.Input())
	}()
}


