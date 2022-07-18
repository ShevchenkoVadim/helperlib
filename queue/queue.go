package queue

import (
	"github.com/ShevchenkoVadim/helperlib/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const delay = 3

type Rabbit struct {
	conn         *amqp.Connection
	ch           *amqp.Channel
	Uri          string
	Queue        string
	WaitChannel  chan bool
	ChanDelivery <-chan amqp.Delivery
}

func (r *Rabbit) writeToWaitChannel() {
	go func() {
		r.WaitChannel <- true
	}()
}

func logWrapper(msg ...any) {
	if config.C.Debug {
		log.Println(msg...)
	}
}

func (r *Rabbit) TestPortRabbitMQ() {
	logWrapper("func TestPortRabbitMQ")
	uri := strings.Split(r.Uri, "@")
	log.Println("!!!!", uri)
	conn, err := net.DialTimeout("tcp", uri[1], time.Second)
	if conn != nil {
		defer conn.Close()
	}
	if err != nil {
		logWrapper(uri[1], "Net err ", err)
		if r.ch != nil && !r.ch.IsClosed() {
			r.ch.Close()
			r.ch = nil
		}
		if r.ch != nil && !r.conn.IsClosed() {
			r.conn.Close()
			r.conn = nil
		}
		connErr := r.Channel()
		for connErr != nil {
			connErr = r.Channel()
		}
		r.writeToWaitChannel()
		logWrapper("Err not nil at TespPortRabbitMQ")
	} else {
		logWrapper("Write to wait at TestPortRabbitMQ")
		r.writeToWaitChannel()
	}
}

func (r *Rabbit) Connect() error {
	logWrapper("Connect to queue server: ", r.Uri)
	if r.conn != nil {
		r.conn.Close()
	}
	conn, err := amqp.Dial(r.Uri)
	if err != nil {
		return err
	}
	logWrapper("Connected to: ", r.Uri)
	r.conn = conn
	go func() {
		for {
			err = <-r.conn.NotifyClose(make(chan *amqp.Error))
			logWrapper("Some error with connection closed at func Connect")
			for {
				time.Sleep(delay * time.Second)
				conn, err := amqp.Dial(r.Uri)
				if err == nil {
					r.conn = conn
					logWrapper("reconnect success")
					break
				}
				logWrapper("reconnect failed at func Connect: ", err)
			}
		}
	}()
	return nil
}

func (r *Rabbit) Channel() error {
	logWrapper("Open queue channel: ", r.Queue)
	err := r.Connect()
	if err != nil {
		return err
	}

	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}
	logWrapper("Channel created: ", r.Queue)
	r.ch = ch
	go func() {
		for {
			err = <-r.ch.NotifyClose(make(chan *amqp.Error))
			logWrapper("Found error at func Channel. From NotifyClose")
			for {
				time.Sleep(delay * time.Second)
				err = r.Connect()
				if err != nil {
					continue
				}
				ch, err := r.conn.Channel()
				if err == nil {
					logWrapper("channel recreate success: ", r.Queue)
					r.ch = ch
					r.writeToWaitChannel()
					break
				}
				logWrapper("channel recreate failed: ", err)
			}
		}
	}()
	return nil
}

func (r *Rabbit) Publish(msg []byte) error {
	logWrapper("Publish message to: ", r.Queue)
	if r.conn == nil || r.ch == nil {
		r.Channel()
	}
	r.TestPortRabbitMQ()
	logWrapper("Waiting while port is opened for publish message to queue: ", r.Queue)
	<-r.WaitChannel
	if r.ch != nil {
		logWrapper("QueuerDeclare: ", r.Queue)
		_, err := r.ch.QueueDeclare(
			r.Queue,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			logWrapper("CHANNEL ERROR at func Publish ", err)
			r.Channel()
			r.Publish(msg)
		} else {
			err = r.ch.Publish(
				"",
				r.Queue,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/json",
					Body:        msg,
				})
			if err != nil {
				logWrapper("PUBLISH ERROR at func Publish", err)
			}
			return err
		}

	} else {
		logWrapper("Channel is nil")
		r.Channel()
		r.Publish(msg)
	}
	return nil
}

func (r *Rabbit) Consume() {
	logWrapper("Consume to queue: ", r.Queue)
	if r.conn == nil || r.ch == nil {
		r.Channel()
	}
	r.TestPortRabbitMQ()
	logWrapper("Waiting while port is opened for consume queue: ", r.Queue)
	<-r.WaitChannel
	if r.ch != nil {
		logWrapper("QueuerDeclare")
		_, err := r.ch.QueueDeclare(
			r.Queue,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			logWrapper("CHANNEL ERROR at func Consume ", err)
			r.Channel()
			r.Consume()
		} else {
			err = r.ch.Qos(
				1,
				0,
				false,
			)

			if err != nil {
				logWrapper("QOS ERROR at func Consume ", err)
				r.Channel()
				r.Consume()
			}
			logWrapper("Consume again to queue: ", r.Queue)
			r.ChanDelivery, err = r.ch.Consume(
				r.Queue,
				filepath.Base(os.Args[0]),
				false,
				false,
				false,
				false,
				nil,
			)
			if err != nil {
				logWrapper("CONSUME ERROR at func Consume ", err)
			}
		}
	} else {
		logWrapper("Channel is nil")
		r.Channel()
		r.Consume()
	}
}
