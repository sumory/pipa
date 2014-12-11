package chainpipe

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
}

func NewPipa() (pl *Pipa) {
	pl = new(Pipa)
	pl.status = WAITING
	pl.pipes = []IPipe{}
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


func (self *Pipa) Connect() (stop chan bool) {
	var output chan interface{}
	output, _ = self.source.ConnectSource()
	for _, p := range self.pipes {
		output, _ = p.ConnectPipe(output)
	}
	stop, _ = self.sink.ConnectSink(output)
	self.status = READY
	return
}


func (self *Pipa) ConnectPipe(input chan interface{}) (output chan interface{}) {
	//implements by subclass
	return
}

func (self *Pipa) Status() int {return self.status}



func (self *Pipa) AddSource(src ISource) error {
	if self.status == RUNNING {
		return errors.New("Abandon 'AddSource' when RUNNING\n")
	}
	if self.source != nil {
		return errors.New("source already added")
	}
	self.source = src
	return nil
}

func (self *Pipa) AddPipe(p IPipe) error {
	if self.status == RUNNING {
		return errors.New("Abandon 'AddPipe' when RUNNING")
	}
	self.pipes = append(self.pipes, p)
	return nil
}

func (self *Pipa) AddSink(sk ISink) error {
	if self.status == RUNNING {
		return errors.New("Abandon 'AddSink' when RUNNING")
	}
	if self.sink != nil {
		return errors.New("sink already added")
	}
	self.sink = sk
	return nil
}


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

func (self *Pipa) Pipe(p IPipe) *Pipa {
	if self.status == RUNNING {
		return self
	}
	self.pipes = append(self.pipes, p)
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
