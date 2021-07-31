package fs

import (
	"bufio"
	flow "go-flow"
	"os"
)

type FileSource struct {
	name string
	in   chan interface{}
}

func NewFileSource(filename string) *FileSource {
	source := &FileSource{
		name: filename,
		in:   make(chan interface{}),
	}
	source.init()
	return source
}

func (fs *FileSource) init() {
	go func() {
		f, err := os.Open(fs.name)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		reader := bufio.NewReader(f)
		for {
			l, isPrefix, err := reader.ReadLine()
			if err != nil {
				close(fs.in)
				break
			}
			var msg string
			if isPrefix {
				msg = string(l)
			} else {
				msg = string(l) + "\n"
			}
			fs.in <- msg
		}
	}()
}

func (fs *FileSource) Output() <-chan interface{} {
	return fs.in
}

func (fs *FileSource) Via(_flow flow.DataFlow) flow.DataFlow {
	flow.FlowTo(fs, _flow)
	return _flow
}

type FileTarget struct {
	filename string
	in       chan interface{}
}

func NewFileTarget(filename string)*FileTarget{
	ft := &FileTarget{
		filename: filename,
		in:       make(chan interface{}),
	}
	go ft.init()
	return ft
}

func (ft *FileTarget)init(){
	go func() {
		f,err := os.OpenFile(ft.filename,os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		for elem := range ft.in{
			_,err:=f.WriteString(elem.(string))
			if err != nil {
				
			}
		}
	}()
}
