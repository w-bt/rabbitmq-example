package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
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
		"ExDirect", // name
		"direct",   // type
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
	queueList := map[string][]string{
		"TestDirect1": []string{"info"},
		"TestDirect2": []string{"warning", "critical"},
	}

	for queue, values := range queueList {
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

		for _, key := range values {
			err = ch.QueueBind(
				q.Name,     // queue name
				key,        // routing key
				"ExDirect", // exchange
				false,
				nil,
			)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	// attempt to publish a message to the queue!
	body := bodyFrom(os.Args)
	err = ch.Publish(
		"ExDirect",
		severityFrom(os.Args), // routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	log.Println(severityFrom(os.Args))
	log.Println(body)

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

func severityFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "info"
	} else {
		s = os.Args[1]
	}
	return s
}
