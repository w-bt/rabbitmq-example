package main

import (
	"fmt"

	"github.com/NeowayLabs/wabbit"
	"github.com/NeowayLabs/wabbit/amqp"
)

func NewConnection() (wabbit.Conn, error) {
	return amqp.Dial("amqp://guest:guest@localhost:5672/")
}

func main() {
	fmt.Println("Go RabbitMQ Tutorial")
	conn, err := NewConnection()
	if err != nil {
		fmt.Println("Failed Initializing Broker Connection")
		panic(err)
	}

	ch, err := initChannel(conn)
	if err != nil {
		fmt.Println("Failed Initializing Channel")
		panic(err)
	}

	err = publish(ch)
	if err != nil {
		fmt.Println("Failed Publishing Message")
		panic(err)
	}
}

func initChannel(conn wabbit.Conn) (wabbit.Channel, error) {
	// Let's start by opening a channel to our RabbitMQ instance
	// over the connection we have already established
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		return &amqp.Channel{}, err
	}

	err = ch.ExchangeDeclare(
		"ExDirect", // name
		"direct",   // type
		nil,        // opt
	)
	if err != nil {
		fmt.Println(err)
		defer ch.Close()
		return ch, err
	}

	// with this channel open, we can then start to interact
	// with the instance and declare Queues that we can publish and
	// subscribe to
	queueList := map[string][]string{
		"TestDirect1": []string{"info"},
	}

	for queue, values := range queueList {
		q, err := ch.QueueDeclare(
			queue,
			wabbit.Option{
				"durable":    true,
				"autoDelete": false,
				"exclusive":  false,
				"noWait":     false,
			},
		)
		// We can print out the status of our Queue here
		// this will information like the amount of messages on
		// the queue
		fmt.Println(q)
		// Handle any errors if we were unable to create the queue
		if err != nil {
			fmt.Println(err)
			ch.Close()
			return ch, err
		}

		for _, key := range values {
			err = ch.QueueBind(
				q.Name(),   // queue name
				key,        // routing key
				"ExDirect", // exchange
				wabbit.Option{
					"noWait": false,
				},
			)
			if err != nil {
				fmt.Println(err)
				ch.Close()
				return ch, err
			}
		}
	}

	return ch, nil
}

func publish(ch wabbit.Channel) error {
	// attempt to publish a message to the queue!
	err := ch.Publish(
		"ExDirect",
		"info", // routing key
		[]byte("ini pesan"),
		wabbit.Option{
			"deliveryMode": 2,
			"contentType":  "text/plain",
		},
	)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Successfully Published Message to Queue")

	return nil
}
