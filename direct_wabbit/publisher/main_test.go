package main

import (
	"testing"

	"github.com/NeowayLabs/wabbit"
	"github.com/NeowayLabs/wabbit/amqptest"
	"github.com/NeowayLabs/wabbit/amqptest/server"
)

func TestPublish(t *testing.T) {
	amqpuri := "amqp://guest:guest@localhost:35672/%2f"

	server := server.NewServer(amqpuri)
	server.Start()
	defer server.Stop()

	conn, err := amqptest.Dial(amqpuri)

	if err != nil || conn == nil {
		t.Error(err)
		return
	}

	// core test 1
	ch, err := initChannel(conn)
	if ch == nil {
		t.Error("channel is empty")
	}
	if err != nil {
		t.Error("init channel failed")
	}
	// end core test 1

	deliveries, err := consume(ch)
	if err != nil {
		t.Error(err)
		return
	}

	// core test 2
	err = publish(ch)
	if ch == nil {
		t.Error("channel is empty")
	}
	// end core test 2

	data := <-deliveries

	if string(data.Body()) != "ini pesan" {
		t.Errorf("Failed to publish message to specified route")
		return
	}
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
