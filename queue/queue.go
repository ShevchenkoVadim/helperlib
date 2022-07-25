package queue

import (
	"github.com/ShevchenkoVadim/helperlib/config"
	"github.com/ShevchenkoVadim/helperlib/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Rabbit struct {
	conn         *amqp.Connection
	ch           *amqp.Channel
	Delay        time.Duration
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
	utils.LogWrapper("func TestPortRabbitMQ")
	uri := strings.Split(r.Uri, "@")
	conn, err := net.DialTimeout("tcp", uri[1], time.Second)
	if conn != nil {
		defer conn.Close()
	}
	if err != nil {
		utils.LogWrapper(uri[1], "Net err ", err)
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
		utils.LogWrapper("Err not nil at TespPortRabbitMQ")
	} else {
		utils.LogWrapper("Write to wait at TestPortRabbitMQ")
		r.writeToWaitChannel()
	}
}

func (r *Rabbit) Connect() error {
	//logWrapper("Connect to queue server: ", r.Uri)
	if r.conn != nil {
		r.conn.Close()
	}
	conn, err := amqp.Dial(r.Uri)
	if err != nil {
		return err
	}
	utils.LogWrapper("Connected to: ", r.Uri)
	r.conn = conn
	go func() {
		for {
			err = <-r.conn.NotifyClose(make(chan *amqp.Error))
			utils.LogWrapper("Some error with connection closed at func Connect")
			for {
				time.Sleep(r.Delay * time.Second)
				utils.LogWrapper("@@@@", r.conn.IsClosed())
				if r.conn.IsClosed() {
					//conn, err := amqp.Dial(r.Uri)
					//if err == nil {
					//	r.conn = conn
					//	logWrapper("reconnect success")
					//	break
					//}
					utils.LogWrapper("reconnect failed at func Connect: ")
					time.Sleep(1 * time.Second)
				} else {
					utils.LogWrapper("reconnect success!!!")
					break
				}
			}
		}
	}()
	return nil
}

func (r *Rabbit) Channel() error {
	utils.LogWrapper("Open queue channel: ", r.Queue)
	err := r.Connect()
	if err != nil {
		return err
	}

	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}
	utils.LogWrapper("Channel created: ", r.Queue)
	r.ch = ch
	go func() {
		for {
			err = <-r.ch.NotifyClose(make(chan *amqp.Error))
			utils.LogWrapper("Found error at func Channel. From NotifyClose")
			for {
				time.Sleep(r.Delay * time.Second)
				err = r.Connect()
				if err != nil {
					continue
				}
				ch, err := r.conn.Channel()
				if err == nil {
					utils.LogWrapper("channel recreate success: ", r.Queue)
					r.ch = ch
					r.writeToWaitChannel()
					break
				}
				utils.LogWrapper("channel recreate failed: ", err)
			}
		}
	}()
	return nil
}

func (r *Rabbit) Publish(msg []byte) error {
	utils.LogWrapper("Publish message to: ", r.Queue)
	if r.conn == nil || r.ch == nil {
		r.Channel()
	}
	r.TestPortRabbitMQ()
	utils.LogWrapper("Waiting while port is opened for publish message to queue: ", r.Queue)
	<-r.WaitChannel
	if r.ch != nil {
		utils.LogWrapper("QueuerDeclare: ", r.Queue)
		_, err := r.ch.QueueDeclare(
			r.Queue,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			utils.LogWrapper("CHANNEL ERROR at func Publish ", err)
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
				utils.LogWrapper("PUBLISH ERROR at func Publish", err)
			}
			return err
		}

	} else {
		utils.LogWrapper("Channel is nil")
		r.Channel()
		r.Publish(msg)
	}
	return nil
}

func (r *Rabbit) Consume() {
	utils.LogWrapper("Consume to queue: ", r.Queue)
	if r.conn == nil || r.ch == nil {
		r.Channel()
	}
	r.TestPortRabbitMQ()
	utils.LogWrapper("Waiting while port is opened for consume queue: ", r.Queue)
	<-r.WaitChannel
	if r.ch != nil {
		utils.LogWrapper("QueuerDeclare")
		_, err := r.ch.QueueDeclare(
			r.Queue,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			utils.LogWrapper("CHANNEL ERROR at func Consume ", err)
			r.Channel()
			r.Consume()
		} else {
			err = r.ch.Qos(
				1,
				0,
				false,
			)

			if err != nil {
				utils.LogWrapper("QOS ERROR at func Consume ", err)
				r.Channel()
				r.Consume()
			}
			utils.LogWrapper("Consume again to queue: ", r.Queue)
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
				utils.LogWrapper("CONSUME ERROR at func Consume ", err)
			}
		}
	} else {
		utils.LogWrapper("Channel is nil")
		r.Channel()
		r.Consume()
	}
}
