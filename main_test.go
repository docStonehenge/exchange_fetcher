// +build integration

package main

import (
	"bytes"
	"context"
	"github.com/streadway/amqp"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestParsingIndicesCorrectlyFromMessageQueue(t *testing.T) {
	testConn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		t.Fatal("Error while trying to set test connection to RabbitMQ")
	}

	defer testConn.Close()

	testChan, err := testConn.Channel()
	if err != nil {
		t.Fatal("Error while trying to set test channel to connection")
	}

	defer testChan.Close()

	pubQueue, err := testChan.QueueDeclare("exchange_fetcher.indices.requests", false, true, false, false, nil)
	if err != nil {
		t.Fatal("Error while trying to set test publishing queue connection", false, true, false, false, nil)
	}

	subQueue, err := testChan.QueueDeclare("exchange_fetcher.indices.results", false, true, false, false, nil)
	if err != nil {
		t.Fatal("Error while trying to set test subscription queue connection")
	}

	testChan.Publish(
		"",
		pubQueue.Name,
		false, false,
		amqp.Publishing{ContentType: "application/json", Body: []byte("{\"indices\":[\"AAPL\"]}")},
	)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "exchange_fetcher", "-mq")
	var input bytes.Buffer
	var output bytes.Buffer

	cmd.Stdout = &output
	cmd.Stdin = &input

	cmd.Run()

	if !strings.Contains(output.String(), "Connected successfully to port 5672") {
		t.Fatal("Error while trying to connect to RabbitMQ server")
	}

	if !strings.Contains(output.String(), "Channel is now opened...") {
		t.Fatal("Error while opening channel")
	}

	if !strings.Contains(output.String(), "Receiving indices on queue 'exchange_fetcher.indices.requests'") {
		t.Fatal("Error on output requests queue message")
	}

	if !strings.Contains(output.String(), "Publishing results on queue 'exchange_fetcher.indices.results'") {
		t.Fatal("Error on output results queue message")
	}

	if !strings.Contains(output.String(), "Waiting for indices.") {
		t.Fatal("Error happened when trying to open subscriber and preparing application")
	}

	msgs, err := testChan.Consume(subQueue.Name, "", true, false, false, false, nil)

	msg := <-msgs

	if !strings.Contains(string(msg.Body), "{\"Apple Inc.\":{\"Name\":\"Apple Inc.\",\"Symbol\":\"AAPL\"") {
		t.Fatalf("Message sent to client should be a JSON with correct values, but is %s", string(msg.Body))
	}
}

func TestParsingIndicesCorrectlyFromCommandLine(t *testing.T) {
	cmd := exec.Command("exchange_fetcher", "-indices", "{\"indices\":[\"AAPL\"]}")
	var output bytes.Buffer

	cmd.Stdout = &output

	if err := cmd.Run(); err == nil {
		if strings.Contains(output.String(), "Connecting to AMQP server...") {
			t.Fatal("Using -c flag should not start MQ connection, but connection has been started")
		}

		if !strings.Contains(output.String(), "{\"Apple Inc.\":{\"Name\":\"Apple Inc.\",\"Symbol\":\"AAPL\"") {
			t.Fatalf("Application should print correct return JSON, but printed only: %s", output.String())
		}
	} else {
		t.Fatal(err)
	}
}
