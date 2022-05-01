package logcollector

import (
	"encoding/json"
	"fmt"
	"github.com/ShevchenkoVadim/helperlib/config"
	"github.com/ShevchenkoVadim/helperlib/queue"
	"github.com/ShevchenkoVadim/helperlib/sfotypes"
	"log"
	"os"
	"path/filepath"
	"time"
)

func sendLogToQueue(service string, msg string) {
	data := &sfotypes.LogMsg{
		ServiceName: service,
		Msg:         msg,
		TimeStamp:   time.Now().Unix(),
	}
	jsonData, err := json.Marshal(&data)
	if err != nil {
		log.Panicln(err)
	} else {
		publisher := queue.Rabbit{Uri: config.C.MQ.Url, Queue: config.C.LogQueue, WaitChannel: make(chan bool)}
		publisher.Publish(jsonData)
	}
}

func SendLog(v ...any) {
	sendLogToQueue(filepath.Base(os.Args[0]), fmt.Sprintln(v...))
}
