package main

import (
	"strings"
	"time"
	"fmt"
)

type LogProcess struct {
	rc          chan string
	wc          chan string
	path        string
	influxDBDsn string
}

func (l *LogProcess) ReadFromFile() {
	//读取数据
	line := "message"
	l.rc <- line
}

func (l *LogProcess) Process() {
	//解析数据
	data := <-l.rc
	l.wc <- strings.ToUpper(data)
}

func (l *LogProcess) WriteToInfluxDB() {
	fmt.Println(<-l.wc)
}

func main() {
	lp := &LogProcess{
		rc:          make(chan string),
		wc:          make(chan string),
		path:        "",
		influxDBDsn: "",
	}

	go lp.ReadFromFile()
	go lp.Process()
	go lp.WriteToInfluxDB()
	time.Sleep(1 * time.Second)
}
