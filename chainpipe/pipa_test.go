package chainpipe

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

// Source
type SampleSource struct {
	msg    string
	size   int
	status int
	input  chan interface{}
	output chan interface{}
}

func NewSampleSource(msg string, size int) (ds *SampleSource) {
	ds = &SampleSource{}
	ds.size = size
	ds.msg = msg
	ds.status = WAITING
	return
}

func (self *SampleSource) Status() int {
	return self.status
}

func (self *SampleSource) Run() {
	go func() {
		for i := 0; i < self.size; i++ {
			self.output <- self.msg
		}
		close(self.output)
	}()
	self.status = RUNNING
}

func (self *SampleSource) ConnectSource() (output chan interface{}, err error) {
	self.output = make(chan interface{})
	self.status = READY
	return self.output, nil
}

// Sink
type SampleSink struct {
	status int
	items  int
	stop   chan bool
	input  chan interface{}
	msg    interface{}
}

func NewSampleSink() (ns *SampleSink) {
	ns = &SampleSink{}
	ns.status = WAITING
	ns.items = 0
	return
}

func (self *SampleSink) Status() int {
	return self.status
}

func (self *SampleSink) Run() {
	go func() {
		for {
			stuff := <-self.input
			if stuff == nil {
				break
			}else {
				self.msg = stuff
			}
			self.items++
		}
		self.stop <- true
		close(self.stop)
	}()
	self.status = RUNNING
}

func (self *SampleSink) ConnectSink(input chan interface{}) (stop chan bool, err error) {
	self.stop = make(chan bool)
	self.input = input
	self.status = READY
	return self.stop, nil
}

// pipe
type SamplePipe struct {

	status int
	output chan interface{}
	input  chan interface{}
}

func NewSamplePipe() (p *SamplePipe) {
	p = &SamplePipe{}
	return
}

func (p *SamplePipe) Status() int {
	return p.status
}

func (p *SamplePipe) Run() {
	go func() {
		for {
			stuff, ok := <-p.input
			if !ok {
				break
			}
			stufs := stuff.(string)
			_ = stufs
			p.output <- stuff
		}
		close(p.output)
	}()
	p.status = RUNNING
}

func (p *SamplePipe) ConnectPipe(input chan interface{}) (output chan interface{}, err error) {
	p.output = make(chan interface{})
	p.input = input
	p.status = READY
	return p.output, nil
}

func TestPipa(t *testing.T) {
	Convey("test pipa", t, func() {
			content := "hello pipa!"
			count := 100

			pl := NewPipa()
			source := NewSampleSource(content, count)
			sink := NewSampleSink()
			pipe1 := NewSamplePipe()
			pipe2 := NewSamplePipe()

			//stop := pl.Source(source).Pipe(pipe1).Pipe(pipe2).Sink(sink).Connect().Run()
			stop := pl.Source(source).Pipe(pipe1,pipe2).Sink(sink).Connect().Run()//all parts Connect and start running
			<-stop //wait for stop

			So(sink.items, ShouldEqual, count)
			So(sink.msg, ShouldEqual, content)
		})
}
