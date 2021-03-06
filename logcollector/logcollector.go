package logcollector

import (
	"encoding/json"
	"fmt"
	"github.com/ShevchenkoVadim/helperlib/queue"
	"github.com/ShevchenkoVadim/helperlib/sfotypes"
	"log"
	"os"
	"path/filepath"
	"time"
)

type LogHelper struct {
	Uri       string
	LogQueue  string
	HostName  string
	publisher queue.Rabbit
	logger    chan string
}

func (l *LogHelper) sendLogToQueue(service string, msg string) {
	data := &sfotypes.LogMsg{
		HostName:    l.HostName,
		ServiceName: service,
		Msg:         msg,
		TimeStamp:   time.Now().Unix(),
	}
	jsonData, err := json.Marshal(&data)
	if err != nil {
		log.Panicln(err)
	} else {
		l.publisher.Publish(jsonData)
	}
}

func (l *LogHelper) InitSendLog(url string) {
	if url != "" {
		l.Uri = url
	}
	l.publisher = queue.Rabbit{Uri: l.Uri, Queue: l.LogQueue, WaitChannel: make(chan bool)}
	l.publisher.Channel()
	l.logger = make(chan string)
	go func() {
		for {
			msg := <-l.logger
			l.sendLogToQueue(filepath.Base(os.Args[0]), msg)
		}
	}()
}

func (l *LogHelper) SendLog(v ...any) {
	go func() {
		l.logger <- fmt.Sprint(v...)
	}()
}
