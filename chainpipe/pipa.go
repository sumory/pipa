package chainpipe

import (
	//"errors"
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

type IPipe interface {
Runable
	ConnectPipe(input chan interface{}) (output chan interface{}, err error)
}

type ISource interface {
Runable
	ConnectSource() (output chan interface{}, err error)
}

type ISink interface {
Runable
	ConnectSink(input chan interface{}) (stop chan bool, err error)
}


type Pipa struct{
	status int
	source ISource
	pipes  []IPipe
	sink   ISink
	all    []Runable
	stop chan bool
}

func NewPipa() (pl *Pipa) {
	pl = new(Pipa)
	pl.status = WAITING
	pl.pipes = []IPipe{}
	return
}

func (self *Pipa) Run() chan bool {
	self.source.Run()
	for _, p := range self.pipes {
		p.Run()
	}
	self.sink.Run()
	self.status = RUNNING
	return self.stop
}


func (self *Pipa) Connect() *Pipa  {
	var output chan interface{}
	output, _ = self.source.ConnectSource()
	for _, p := range self.pipes {
		output, _ = p.ConnectPipe(output)
	}
	self.stop, _ = self.sink.ConnectSink(output)
	self.status = READY
	return self
}


func (self *Pipa) ConnectPipe(input chan interface{}) (output chan interface{}) {
	//implements by subclass
	return
}

func (self *Pipa) Status() int {return self.status}


func (self *Pipa) Source(src ISource) *Pipa {
	if self.status == RUNNING {
		return self
	}
	if self.source != nil {
		return self
	}
	self.source = src
	return self
}

func (self *Pipa) Pipe(p... IPipe) *Pipa {
	if self.status == RUNNING {
		return self
	}
	self.pipes = append(self.pipes, p...)
	return self
}

func (self *Pipa) Sink(sk ISink) *Pipa {
	if self.status == RUNNING {
		return self
	}
	if self.sink != nil {
		return self
	}
	self.sink = sk
	return self
}
