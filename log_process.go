package main

import (
	"strings"
	"time"
	"fmt"
	"os"
	"bufio"
	"io"
	"log"
	"regexp"
	"strconv"
	"net/url"
)

type Reader interface {
	Read(rc chan []byte)
}

type Writer interface {
	Write(rc chan *Message)
}

type ReadFromFile struct {
	path string
}

type WriteToInfluxDB struct {
	influxDBDsn string
}

type LogProcess struct {
	rc    chan []byte
	wc    chan *Message
	read  Reader
	write Writer
}

type Message struct {
	Host string

	TimeLocal                  time.Time
	Method, Resource, Protocol string
	Status                     string
	BytesSent                  int
	Scheme                     string
	Url                        string
	//UpstreamTime, RequestTime  float64
}

/*14.215.176.15 - - [23/APR/2018:21:43:39 +0800]
"GET /CONSOLE/DIST/LAYOUTS/DEFAULT.F5E6C608DE637CAC3F50.JS HTTP/1.1"
200 5665 "HTTPS://WWW.XHXBLOG.CN/?B3ID=H9OXZSYM"
"MOZILLA/5.0 (WINDOWS NT 6.1; WOW64; RV:43.0) GECKO/20100101 FIREFOX/43.0" "-"*/

func (r *ReadFromFile) Read(rc chan []byte) {
	//读取数据

	f, err := os.Open(r.path)
	if err != nil {
		panic(fmt.Sprintf("open file fail:%s", err.Error()))
	}

	//从文件末尾逐行读取文件内容当在空文件中需要加入此注释掉的代码
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

	r := regexp.MustCompile(`([\d\.]+)\s+([^\[]+)\s+([^\[]+)\s+\[([^\]]+)\]\s+\"([^"]+)\"\s+(\d{3})\s+(\d+)\s+\"([^"]+)\"\s+\"([^"]+)\"\s+`)

	/**
	  测试flag
	*/
	//flag := 0
	loc, _ := time.LoadLocation("Asia/Shanghai")
	for v := range l.rc {

		ret := r.FindStringSubmatch(string(v))

		/*测试**/
		/*
		if flag != 2 {
			sp := strings.Split(ret[5], " ")

			fmt.Println(ret)
			fmt.Println("Host:", ret[1])

			fmt.Println("LocalTime:", ret[4])

			fmt.Println("Method:", sp[0])
			uu, _ := url.Parse(sp[1])
			fmt.Println(uu.Path)
			fmt.Println("Path:", sp[1])
			fmt.Println("Protocol:", sp[2])
			fmt.Println("Status:", ret[6])
			fmt.Println("BytesSent:", ret[7])
			fmt.Println("Scheme:", ret[9])
			fmt.Println(ret[8])
			flag++
		}
		*/
		if len(ret) != 10 {

			log.Println("FindStringSubmatch fail:", string(v))
			continue
		}

		t, err := time.ParseInLocation("02/Jan/2006:15:04:05 +0800", ret[4], loc)
		if err != nil {
			log.Println("ParseInLocation fail:", err.Error(), ret[4])
		}
		message := &Message{}

		//14.215.176.15
		message.Host = ret[1]

		//23/APR/2018:21:43:39 +0800
		message.TimeLocal = t
		//MOZILLA/5.0 (WINDOWS NT 6.1; WOW64; RV:43.0) GECKO/20100101 FIREFOX/43.0
		message.Scheme = ret[9]

		//GET /CONSOLE/DIST/LAYOUTS/DEFAULT.F5E6C608DE637CAC3F50.JS HTTP/1.1
		sp := strings.Split(ret[5], " ")

		if len(sp) != 3 {
			log.Println("strings.Split fail:", ret[5])
			continue
		}
		//请求方法
		message.Method = sp[0]

		//请求路径
		u, err := url.Parse(sp[1])

		if err != nil {
			log.Println("url parse fail:", err)
			continue
		}

		message.Resource = u.Path

		//请求协议
		message.Protocol = sp[2]
		//200
		message.Status = ret[7]

		//HTTPS://WWW.XHXBLOG.CN/?B3ID=H9OXZSYM
		message.Url = ret[8]
		//5665
		message.BytesSent, _ = strconv.Atoi(ret[7])
		l.wc <- message
	}

}

func (w *WriteToInfluxDB) Write(wc chan *Message) {
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
		wc:    make(chan *Message),
		read:  r,
		write: w,
	}

	go lp.read.Read(lp.rc)
	go lp.Process()
	go lp.write.Write(lp.wc)
	time.Sleep(30 * time.Second)
}
