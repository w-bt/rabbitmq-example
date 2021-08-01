package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"os"
	"strings"
)

func main() {
	fmt.Println("Go RabbitMQ Tutorial")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println("Failed Initializing Broker Connection")
		panic(err)
	}

	// Let's start by opening a channel to our RabbitMQ instance
	// over the connection we have already established
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"ExFanout", // name
		"fanout",   // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		fmt.Println(err)
	}

	// with this channel open, we can then start to interact
	// with the instance and declare Queues that we can publish and
	// subscribe to
	queueList := []string{"TestFanout1", "TestFanout2", "TestFanout3"}

	for _, queue := range queueList {
		q, err := ch.QueueDeclare(
			queue,
			false,
			false,
			false,
			false,
			nil,
		)
		// We can print out the status of our Queue here
		// this will information like the amount of messages on
		// the queue
		fmt.Println(q)
		// Handle any errors if we were unable to create the queue
		if err != nil {
			fmt.Println(err)
		}

		err = ch.QueueBind(
			q.Name,     // queue name
			"",         // routing key
			"ExFanout", // exchange
			false,
			nil,
		)
		if err != nil {
			fmt.Println(err)
		}
	}

	// attempt to publish a message to the queue!
	body := bodyFrom(os.Args)
	err = ch.Publish(
		"ExFanout",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Published Message to Queue")
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}
