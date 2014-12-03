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
	pipe  Pipe
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
	self.pipe.Run()
	self.sink.Run()

	self.status = RUNNING
}


func (self *Pipa) Connect() (stop chan bool) {
	var output chan interface{}
	output, _ = self.source.ConnectSource()
	output, _ = self.pipe.ConnectPipe(output)
	stop, _ = self.sink.ConnectSink(output)
	self.status = READY
	return
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

func (self *Pipa) AddPipe(p Pipe) error {
	if self.status == RUNNING {
		return errors.New("Abandon 'AddPipe' when RUNNING")
	}
	if self.pipe != nil {
		return errors.New("pipe already added")
	}
	self.pipe = p
	return nil
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

