package main

import (
	"fmt"
	"log"

	"github.com/NeowayLabs/wabbit"
	"github.com/NeowayLabs/wabbit/amqp"
)

type Consumer struct {
	conn    wabbit.Conn
	channel wabbit.Channel
	tag     string
}

func main() {
	c := &Consumer{
		conn:    nil,
		channel: nil,
		tag:     "",
	}

	fmt.Println("Go RabbitMQ Tutorial")
	var err error
	c.conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println("Failed Initializing Broker Connection")
		panic(err)
	}

	c.channel, err = c.conn.Channel()
	if err != nil {
		fmt.Println(err)
	}
	defer c.channel.Close()

	forever := make(chan bool)
	deliveries, err := consume(c.channel)
	if err != nil {
		fmt.Println(err)
		return
	}
	handle(deliveries)
	<-forever
}

func consume(ch wabbit.Channel) (<-chan wabbit.Delivery, error) {
	return ch.Consume(
		"TestDirect1",
		"", // empty consumer tag
		wabbit.Option{
			"noAck":     false,
			"exclusive": false,
			"noLocal":   false,
			"noWait":    false,
		},
	)
}

func handle(deliveries <-chan wabbit.Delivery) {
	for d := range deliveries {
		log.Printf(
			"got %dB delivery: [%v] %q",
			len(d.Body()),
			d.DeliveryTag(),
			d.Body(),
		)
		d.Ack(false)
	}
	log.Printf("handle: deliveries channel closed")
}
