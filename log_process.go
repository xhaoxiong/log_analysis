package main

import (
	"strings"
	"time"
	"fmt"
	"os"
	"bufio"
	"io"
	"log"
)

type Reader interface {
	Read(rc chan []byte)
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
	rc    chan []byte
	wc    chan string
	read  Reader
	write Writer
}

func (r *ReadFromFile) Read(rc chan []byte) {
	//读取数据

	f, err := os.Open(r.path)
	if err != nil {
		panic(fmt.Sprintf("open file fail:%s", err.Error()))
	}

	//从文件末尾逐行读取文件内容
	//f.Seek(0, 2)

	rd := bufio.NewReader(f)

	for {

		line, err := rd.ReadBytes('\n')

		if err == io.EOF {
			log.Println(err)
			time.Sleep(500 * time.Millisecond)
			continue
		} else if err != nil {
			panic(fmt.Sprintf("ReadBytes error: %s", err.Error()))
		}
		rc <- line[:len(line)-1]

	}

}

func (l *LogProcess) Process() {
	//解析数据
	for v := range l.rc {
		l.wc <- strings.ToUpper(string(v))
	}

}

func (w *WriteToInfluxDB) Write(wc chan string) {
	for v := range wc {
		fmt.Println(v)
	}

}

func (l *LogProcess) ReadFromFile() {

}

func (l *LogProcess) WriteToInfluxDB() {

}

func main() {

	r := &ReadFromFile{
		path: "/var/log/nginx/access.log",
	}

	w := &WriteToInfluxDB{
		influxDBDsn: "",
	}

	lp := &LogProcess{
		rc:    make(chan []byte),
		wc:    make(chan string),
		read:  r,
		write: w,
	}

	go lp.read.Read(lp.rc)
	go lp.Process()
	go lp.write.Write(lp.wc)
	time.Sleep(30 * time.Second)
}
