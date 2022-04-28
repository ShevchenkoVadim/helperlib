package queue

import (
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
	conn        *amqp.Connection
	ch          *amqp.Channel
	Uri         string
	Queue       string
	WaitChannel chan bool
}

func (r *Rabbit) writeToWaitChannel() {
	go func() {
		r.WaitChannel <- true
	}()
}

func (r *Rabbit) TestPortRabbitMQ() {
	if r.Uri != "" && len(r.Uri) > 0 {
		uri := strings.Split(r.Uri, "@")
		timeout := time.Second
		conn, err := net.DialTimeout("tcp", uri[1], timeout)
		if conn != nil {
			defer conn.Close()
		}
		if err != nil {
			if r.ch != nil && !r.ch.IsClosed() {
				r.ch.Close()
				r.ch = nil
			}
			if !r.conn.IsClosed() && r.ch != nil {
				r.conn.Close()
				r.conn = nil
			}
			r.Channel()
		}
		r.writeToWaitChannel()
	}
}

func (r *Rabbit) Connect() error {
	conn, err := amqp.Dial(r.Uri)
	if err != nil {
		return err
	}
	r.conn = conn
	go func() {
		for {
			err = <-r.conn.NotifyClose(make(chan *amqp.Error))

			for {
				time.Sleep(delay * time.Second)

				conn, err := amqp.Dial(r.Uri)
				if err == nil {
					r.conn = conn
					log.Println("reconnect success")
					break
				}
				log.Println("reconnect failed ", err)
			}
		}
	}()
	return nil
}

func (r *Rabbit) Channel() error {
	if r.conn != nil {
		r.conn.Close()
	}
	err := r.Connect()
	if err != nil {
		return err
	}

	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}

	r.ch = ch
	go func() {
		for {
			err = <-r.ch.NotifyClose(make(chan *amqp.Error))

			for {
				time.Sleep(delay * time.Second)

				ch, err := r.conn.Channel()
				if err == nil {
					log.Println("channel recreate success")
					r.ch = ch
					r.writeToWaitChannel()
					break
				}
				log.Println("channel recreate failed ", err)
			}
		}
	}()
	return nil
}

func (r *Rabbit) Publish(msg []byte) error {
	if r.conn == nil || r.ch == nil {
		r.Channel()
	}
	r.TestPortRabbitMQ()
	q, err := r.ch.QueueDeclare(
		r.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("CHANNEL ERROR ", err)
	}
	<-r.WaitChannel
	err = r.ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
		})
	if err != nil {
		log.Println("PUBLISH ERROR ", err)
	}

	return err
}

func (r *Rabbit) Consume() (<-chan amqp.Delivery, error) {
	if r.conn == nil || r.ch == nil {
		r.Channel()
	}
	r.TestPortRabbitMQ()

	q, err := r.ch.QueueDeclare(
		r.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("CHANNEL ERROR ", err)
		return nil, err
	}
	<-r.WaitChannel

	err = r.ch.Qos(
		1,
		0,
		false,
	)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	msgs, err := r.ch.Consume(
		q.Name,
		filepath.Base(os.Args[0]),
		false,
		false,
		false,
		false,
		nil,
	)
	return msgs, nil
}
