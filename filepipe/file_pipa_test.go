package filepipe

import (
	"crypto/md5"
	"io"
	"os"
	"bufio"
	"fmt"
	"time"
	"io/ioutil"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func IsBigEndian() bool {
	var i int32 = 0x12345678
	var b byte = byte(i)
	if b == 0x12 {
		return true
	}
	return false
}

func Now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func getFileInfo(filename string) Header {
	fi, err := os.Lstat(filename)
	if err != nil {
		//fmt.Println("info ERROR", err)
		panic(err)
	}

	fileHandle, err := os.Open(filename)
	defer fileHandle.Close()
	if err != nil {
		//fmt.Println("open ERROR", err)
		panic(err)
	}

	h := md5.New()
	_, err = io.Copy(h, fileHandle)

	fileInfo := Header {
		fName : fi.Name(),
		fSize : fi.Size(),
		fPerm : fi.Mode().Perm(),
		fMtime: fi.ModTime(),
		fType : fi.IsDir(),
		fMd5  : fmt.Sprintf("%x", h.Sum(nil)),
	}
	return fileInfo
}

type Header struct{//文件头信息
	fName       string
	fSize       int64
	fMtime      time.Time
	fPerm       os.FileMode
	fMd5        string
	fType       bool
}

type Body struct{
	content []byte
}

type Packet struct{
	id      int32
	Header
	Body
}

// Source
type FileSource struct {
	filename string
	size     int64
	reader *bufio.Reader
	status   int
	input    chan interface{}
	output   chan interface{}
}

func NewFileSource(filename string) (*FileSource) {
	source := &FileSource{}
	source.status = WAITING
	source.filename = filename
	return source
}

func (self *FileSource) Status() int {
	return self.status
}

func (self *FileSource) Run() {
	go func() {
		f, err := os.OpenFile(self.filename, os.O_RDONLY, 0660)
		defer f.Close()
		if err != nil {
			panic(err)
		}
		var n int64
		if fi, err := f.Stat(); err == nil {
			if size := fi.Size(); size < 1e9 {
				n = size
			}else {
				//should return error
			}
		}
		self.size = n
		self.reader = bufio.NewReader(f)
		self.status = RUNNING
		var count int32 = 0
		for {
			buf := make([]byte, 1024)
			count++
			m, err := self.reader.Read(buf)
			if err != nil && err != io.EOF {panic(err)}
			//fmt.Println("read m -->", m)
			if 0 == m {
				self.output <- nil
				break
			} else {
				packet := &Packet{}
				if count == 1 {//第一个包，发送头信息
					packet.Header = getFileInfo(self.filename)
				}
				packet.id = count
				packet.Body.content = buf[:m]
				//fmt.Println("read:", packet)
				self.output <- packet
			}
		}
		close(self.output)
	}()
}

func (self *FileSource) ConnectSource() (output chan interface{}, err error) {
	self.output = make(chan interface{},10)
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

func NewFileSink() *FileSink {
	sink := &FileSink{}
	sink.status = WAITING
	return sink
}

func (self *FileSink) Status() int {
	return self.status
}

func (self *FileSink) Run() {
	var f *os.File
	go func() {
		for {
			content := <-self.input

			if content == nil {
				break
			} else {//write to file
				packet := content.(*Packet)
				//fmt.Println("sink:", packet.Body)
				if packet.id == 1 {//第一个包，创建文件
					self.writer,f= createSinkFile(packet.Header.fName+".sink")
				}
				self.writer.Write(packet.Body.content)
			}
		}
		self.writer.Flush()
		f.Close()
		self.stop <- true
		close(self.stop)
	}()
	self.status = RUNNING
}

func createSinkFile(filename string)( *bufio.Writer ,*os.File){
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0660)
	//defer f.Close()
	if err != nil {
		panic(err)
	}
	writer := bufio.NewWriter(f)
	return writer,f
}

func(self *FileSink) ConnectSink(input chan interface{}) (stop chan bool, err error) {
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

func (p *FilePipe) Run() {
	go func() {
		for {
			stuff, ok := <-p.input
			//fmt.Println("pile:", stuff)
			if !ok {
				break
			}
			p.output <- stuff
		}
		close(p.output)
	}()
	p.status = RUNNING
}

func (p *FilePipe) ConnectPipe(input chan interface{}) (output chan interface{}, err error) {
	p.output = make(chan interface{}, 10)
	p.input = input
	p.status = READY
	return p.output, nil
}

func TestGetFileInfo(t *testing.T) {
	Convey("test GetFileInfo", t, func() {
			So(func(){getFileInfo("10mb.file.not.exist")}, ShouldPanic )
			So("f1c9645dbc14efddc7d8a322685f26eb", ShouldEqual, getFileInfo("10mb.file").fMd5)
			So("10mb.file", ShouldEqual, getFileInfo("10mb.file").fName)
			So(10485760, ShouldEqual, getFileInfo("10mb.file").fSize)
		})
}

func TestTextFilePipa(t *testing.T) {
	Convey("test text file", t, func() {
			pp := NewPipa()
			source := NewFileSource("test.file")
			pp.AddSource(source)

			sink := NewFileSink()
			pp.AddSink(sink)

			pp.AddPipe(NewFilePipe())

			stop := pp.Connect() //all 3 parts Connect*
			pp.Run()             // all 3 parts start
			<-stop
			So("test.file.sink", ShouldEqual, getFileInfo("test.file.sink").fName)
			So(getFileInfo("test.file").fMd5, ShouldEqual, getFileInfo("test.file.sink").fMd5)
		})
}


func TestFileSource(t *testing.T) {
	Convey("test file source", t, func() {
			source := NewFileSource("test.file")
			outchan, err := source.ConnectSource()
			if err != nil {
				t.Log("connect source error", err)
			}else {
				source.Run()
			}

			var result []byte
			for e := range outchan {
				if e == nil {break}
				tmp := e.(*Packet).Body.content
				result = append(result, tmp[:]...)

			}

			f,err:=os.OpenFile("test.file", os.O_RDONLY, 0660)
			defer f.Close()
			c,err:=ioutil.ReadAll(f)
			So(string(result),ShouldEqual,string(c))
			//t.Log("文件内容为：", string(result))
		})
}


func TestBinaryFilePipa(t *testing.T) {
	Convey("test binary file", t, func() {
			pp := NewPipa()
			source := NewFileSource("10mb.file")
			pp.AddSource(source)

			sink := NewFileSink()
			pp.AddSink(sink)

			pp.AddPipe(NewFilePipe())

			stop := pp.Connect() //all 3 parts Connect*
			pp.Run()             // all 3 parts start
			<-stop
			So("10mb.file.sink", ShouldEqual, getFileInfo("10mb.file.sink").fName)
			So(getFileInfo("10mb.file").fMd5, ShouldEqual, getFileInfo("10mb.file.sink").fMd5)
		})
}

