package filepipe

import (
	. "github.com/smartystreets/goconvey/convey"

	"io"
	"os"
	"bufio"
"fmt"
	"testing"
)

// Source
type FileSource struct {
	filename string
	size   int64
	reader *bufio.Reader
	status int
	input  chan interface{}
	output chan interface{}
}

func NewFileSource(filename string)  (*FileSource){
	source := &FileSource{}
	source.status = WAITING
	source.filename = filename
	return source
}

func (self *FileSource) Status() int {
	return self.status
}

func (self *FileSource) Run() error{
	fmt.Println("------------------>",self.filename)
	f, err := os.OpenFile(self.filename, os.O_RDONLY, 0660)
	defer f.Close()
	if err != nil {
		return  err
	}
	var n int64
	if fi, err := f.Stat(); err == nil {
		if size := fi.Size(); size < 1e9 {
			n = size
		}else{
			//should return error
		}
	}
	self.size = n
	reader := bufio.NewReader(f)
	self.status = RUNNING
	buf := make([]byte, 1024)
	go func() {

		for  {
			n,err := reader.Read(buf)
			if err != nil && err != io.EOF {panic(err)}
			fmt.Println("n -->", n)
			if 0==n{
				self.output <- nil
				break
			} else{
				self.output <- buf[:n]
			}
		}
		close(self.output)
	}()

	return nil
}

func (self *FileSource) ConnectSource() (output chan interface{}, err error) {
	self.output = make(chan interface{})
	self.status = READY
	return self.output, nil
}

// Sink
type FileSink struct {
	writer *bufio.Writer
	status int
	stop   chan bool
	input  chan interface{}
}

func NewFileSink()  *FileSink {
	sink := &FileSink{}
	sink.status = WAITING
	return sink
}

func (self *FileSink) Status() int {
	return self.status
}

func (self *FileSink) Run() error{
	f, err :=os.OpenFile("10mb.file.sink", os.O_WRONLY, 0660)
	defer f.Close()
	if err != nil {
		return err
	}

	self.writer = bufio.NewWriter(f)
	go func() {
		for {
			content := <-self.input
			if content == nil {
				break
			} else {//write to file
				self.writer.Write(content.([]byte))
			}
		}
		self.stop <- true
		close(self.stop)
	}()
	self.status = RUNNING
	return nil
}

func (self *FileSink) ConnectSink(input chan interface{}) (stop chan bool, err error) {
	self.stop = make(chan bool)
	self.input = input
	self.status = READY
	return self.stop, nil
}

// pipe
type FilePipe struct {
	status int
	output chan interface{}
	input  chan interface{}
}

func NewFilePipe() (p *FilePipe) {
	p = &FilePipe{}
	return
}

func (p *FilePipe) Status() int {
	return p.status
}

func (p *FilePipe) Run() error{
	go func() {
		for {
			stuff, ok := <-p.input
			if !ok {
				break
			}
			p.output <- stuff
		}
		close(p.output)
	}()
	p.status = RUNNING
	return nil
}

func (p *FilePipe) ConnectPipe(input chan interface{}) (output chan interface{}, err error) {
	p.output = make(chan interface{})
	p.input = input
	p.status = READY
	return p.output, nil
}

func TestFilePipa(t *testing.T) {
	Convey("test file pipa", t, func() {
		pp := NewPipa()

		source := NewFileSource("/data/go/src/github.com/sumory/pipa/filepipe/10mb.file")
		pp.AddSource(source)

		sink := NewFileSink()
		pp.AddSink(sink)

		pp.AddPipe(NewFilePipe())

		stop := pp.Connect() //all 3 parts Connect*
		pp.Run()             // all 3 parts start
		<-stop


	})
}
