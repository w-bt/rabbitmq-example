package main

import (
	"testing"

	"github.com/NeowayLabs/wabbit/amqptest/server"
)

func TestConsume(t *testing.T) {
	vh := server.NewVHost("/")

	ch := server.NewChannel(vh)

	err := ch.ExchangeDeclare("ExDirect", "direct", nil)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = ch.QueueDeclare("TestDirect1", nil)
	if err != nil {
		t.Error(err)
		return
	}

	err = ch.QueueBind("TestDirect1", "info", "ExDirect", nil)
	if err != nil {
		t.Error(err)
		return
	}

	// core test
	deliveries, err := consume(ch)
	if err != nil {
		t.Error(err)
		return
	}
	// end core test

	err = ch.Publish("ExDirect", "info", []byte("ini test"), nil)

	if err != nil {
		t.Error(err)
		return
	}

	data := <-deliveries

	if string(data.Body()) != "ini test" {
		t.Errorf("Failed to publish message to specified route")
		return
	}
}
