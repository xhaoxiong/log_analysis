package main

import (
	"strings"
	"time"
	"fmt"
)

type Reader interface {
	Read(rc chan string)
}

type Writer interface {
	Write(rc chan string)
}

type ReadFromFile struct {
	path string
}

type WriteToInfluxDB struct {
	influxDBDsn string
}

type LogProcess struct {
	rc    chan string
	wc    chan string
	read  Reader
	write Writer
}

func (r *ReadFromFile) Read(rc chan string) {
	//读取数据
	line := "message"
	rc <- line
}

func (l *LogProcess) Process() {
	//解析数据
	data := <-l.rc
	l.wc <- strings.ToUpper(data)
}

func (w *WriteToInfluxDB) Write(wc chan string) {
	fmt.Println(<-wc)
}

func (l *LogProcess) ReadFromFile() {

}

func (l *LogProcess) WriteToInfluxDB() {

}

func main() {

	r := &ReadFromFile{
		path: "",
	}

	w := &WriteToInfluxDB{
		influxDBDsn: "",
	}

	lp := &LogProcess{
		rc:    make(chan string),
		wc:    make(chan string),
		read:  r,
		write: w,
	}


	go lp.read.Read(lp.rc)
	go lp.Process()
	go lp.write.Write(lp.wc)
	time.Sleep(1 * time.Second)
}
