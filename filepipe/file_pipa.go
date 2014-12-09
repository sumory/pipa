package filepipe

import (
	"errors"
	//"fmt"
)

const (
	WAITING = -1
	READY   = 0
	RUNNING = 1
)

type Runable interface {
	Run()
	Status() int
}

type Pipe interface {
Runable
	ConnectPipe(input chan interface{}) (output chan interface{}, err error)
}

type Source interface {
Runable
	ConnectSource() (output chan interface{}, err error)
}

type Sink interface {
Runable
	ConnectSink(input chan interface{}) (stop chan bool, err error)
}


type Pipa struct{
	status int
	source Source
	pipes  []Pipe
	sink   Sink
	all    []Runable
}

func NewPipa() (pl *Pipa) {
	pl = new(Pipa)
	pl.status = WAITING
	return
}

func (self *Pipa) Run() {
	self.source.Run()
	for _, p := range self.pipes {
		p.Run()
	}
	self.sink.Run()

	self.status = RUNNING
}


func (self *Pipa) Connect() (chan bool) {
	var source_output, pipe_output chan interface{}
	source_output, _ = self.source.ConnectSource()

	pipe_output = source_output
	for _, p := range self.pipes {//多管道连接，rafactor later
		pipe_output, _ = p.ConnectPipe(pipe_output)
	}
	stop, _ := self.sink.ConnectSink(pipe_output)
	self.status = READY
	return stop
}


func (self *Pipa) Status() int {return self.status}



func (self *Pipa) AddSource(src Source) error {
	if self.status == RUNNING {
		return errors.New("Abandon 'AddSource' when RUNNING\n")
	}
	if self.source != nil {
		return errors.New("source already added")
	}
	self.source = src
	return nil
}

func (self *Pipa) AddPipe(p Pipe) (error,*Pipa) {
	//if self.status == RUNNING {
		return errors.New("Abandon 'AddPipe' when RUNNING"),nil
	//}
	//self.pipes = append(self.pipes, p)
	//return nil, self
}

func (self *Pipa) AddSink(sk Sink) error {
	if self.status == RUNNING {
		return errors.New("Abandon 'AddSink' when RUNNING")
	}
	if self.sink != nil {
		return errors.New("sink already added")
	}
	self.sink = sk
	return nil
}

