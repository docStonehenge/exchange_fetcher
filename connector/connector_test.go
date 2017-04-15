package connector

import (
	"fmt"
	"github.com/docStonehenge/exchange_fetcher/exchange"
	"github.com/streadway/amqp"
	"os"
	"regexp"
	"strings"
	"testing"
)

var connErrMatch = regexp.MustCompile("There was a problem when opening connection to AMQP")
var channelErrMatch = regexp.MustCompile("There was a problem when opening channel on connection")

func TestOpenConnectionWithSuccess(t *testing.T) {
	os.Setenv("AMQP_USERNAME", "guest")
	os.Setenv("AMQP_PASSWORD", "guest")
	os.Setenv("AMQP_DEFAULT_PORT", "5672")

	connection, err := OpenConnection()

	if connection != nil {
		defer connection.Close()
	} else {
		t.Fatal(
			"There was a problem on opening connection for testing. Maybe RabbitMQ server is down ?! This is an integration test, so it's necessary to have an opened connection to AMQP server.",
		)
	}

	if err != nil {
		t.Fatalf("OpenConnection() should return correct connection to amqp, but returned error: %v", err)
	}

	os.Setenv("AMQP_USERNAME", "")
	os.Setenv("AMQP_PASSWORD", "")
	os.Setenv("AMQP_DEFAULT_PORT", "")
}

func TestOpenConnectionWithError(t *testing.T) {
	_, err := OpenConnection()

	if err == nil {
		t.Fatalf("OpenConnection() should return an error, but error is null: %v", err)
	}

	msg := fmt.Sprintf("%v", err)

	if !connErrMatch.MatchString(msg) {
		t.Fatalf("Error message should be for connection error, but is %s", msg)
	}
}

func TestConnectionError(t *testing.T) {
	_, err := OpenConnection()

	connError := ConnectionError{err: err}
	msg := connError.Error()

	if !connErrMatch.MatchString(msg) {
		t.Fatalf("Error message should have prefix of 'There was a problem when opening connection to AMQP', but message is: %s", msg)
	}
}

func TestOpenChannelWithSuccess(t *testing.T) {
	os.Setenv("AMQP_USERNAME", "guest")
	os.Setenv("AMQP_PASSWORD", "guest")
	os.Setenv("AMQP_DEFAULT_PORT", "5672")

	connection, err := OpenConnection()

	if connection != nil {
		defer connection.Close()
	} else {
		t.Fatal(
			"There was a problem on opening connection for testing. Maybe RabbitMQ server is down ?! This is an integration test, so it's necessary to have an opened connection to AMQP server.",
		)
	}

	channel, err := OpenChannel(connection)
	defer channel.Close()

	if err != nil {
		t.Fatalf("OpenChannel() should return an opened channel, but returned error: %v", err)
	}

	os.Setenv("AMQP_USERNAME", "")
	os.Setenv("AMQP_PASSWORD", "")
	os.Setenv("AMQP_DEFAULT_PORT", "")
}

func TestOpenChannelWithError(t *testing.T) {
	os.Setenv("AMQP_USERNAME", "guest")
	os.Setenv("AMQP_PASSWORD", "guest")
	os.Setenv("AMQP_DEFAULT_PORT", "5672")

	connection, err := OpenConnection()

	if connection != nil {
		connection.Close()
	} else {
		t.Fatal(
			"There was a problem on opening connection for testing. Maybe RabbitMQ server is down ?! This is an integration test, so it's necessary to have an opened connection to AMQP server.",
		)
	}

	_, err = OpenChannel(connection)

	if err == nil {
		t.Fatalf("OpenChannel() should return an error, but returned nothing: %v", err)
	}

	os.Setenv("AMQP_USERNAME", "")
	os.Setenv("AMQP_PASSWORD", "")
	os.Setenv("AMQP_DEFAULT_PORT", "")
}

func TestChannelError(t *testing.T) {
	os.Setenv("AMQP_USERNAME", "guest")
	os.Setenv("AMQP_PASSWORD", "guest")
	os.Setenv("AMQP_DEFAULT_PORT", "5672")

	connection, err := OpenConnection()

	if connection != nil {
		connection.Close()
	} else {
		t.Fatal(
			"There was a problem on opening connection for testing. Maybe RabbitMQ server is down ?! This is an integration test, so it's necessary to have an opened connection to AMQP server.",
		)
	}

	_, err = OpenChannel(connection)
	channelError := &ChannelError{err: err}

	msg := channelError.Error()

	if !channelErrMatch.MatchString(msg) {
		t.Fatalf("Channel error message should start with 'There was a problem when opening channel on connection', but is %s", msg)
	}

	os.Setenv("AMQP_USERNAME", "")
	os.Setenv("AMQP_PASSWORD", "")
	os.Setenv("AMQP_DEFAULT_PORT", "")
}

func TestDefineQueue(t *testing.T) {
	os.Setenv("AMQP_USERNAME", "guest")
	os.Setenv("AMQP_PASSWORD", "guest")
	os.Setenv("AMQP_DEFAULT_PORT", "5672")

	connection, err := OpenConnection()

	if connection != nil {
		defer connection.Close()
	} else {
		t.Fatal(
			"There was a problem on opening connection for testing. Maybe RabbitMQ server is down ?! This is an integration test, so it's necessary to have an opened connection to AMQP server.",
		)
	}

	channel, err := OpenChannel(connection)

	if channel != nil {
		defer channel.Close()
	} else {
		t.Fatalf("There was a problem on opening channel: %v", err)
	}

	queue, err := DefineQueue(channel, "test_queue")

	if err != nil {
		t.Fatalf("DefineQueue() should return a queue, but returned error: %v", err)
	}

	if queue.Name != "test_queue" {
		t.Fatalf("Queue name should be %s, but is %s", "test_queue", queue.Name)
	}

	if _, err := channel.QueueDelete("test_queue", false, false, false); err != nil {
		t.Fatal("Queue 'test_queue' should be deleted.")
	}

	os.Setenv("AMQP_USERNAME", "")
	os.Setenv("AMQP_PASSWORD", "")
	os.Setenv("AMQP_DEFAULT_PORT", "")
}

// These are integration tests that must be run only on controlled environment, since amqp connection cannot be closed from within the test.
// To run these tests:
//   * uncomment the imports
//   * open RabbitMQ dashboard at http://localhost:15672 and force-close connection after test run, to exit it.

func TestOpenSubscriber(t *testing.T) {
	os.Setenv("AMQP_USERNAME", "guest")
	os.Setenv("AMQP_PASSWORD", "guest")
	os.Setenv("AMQP_DEFAULT_PORT", "5672")

	connection, err := OpenConnection()
	defer connection.Close()

	if err != nil {
		t.Fatal(
			"There was a problem on opening connection for testing. Maybe RabbitMQ server is down ?! This is an integration test, so it's necessary to have an opened connection to AMQP server.",
		)
	}

	channel, err := OpenChannel(connection)
	defer channel.Close()

	if err != nil {
		t.Fatalf("There was a problem on opening channel: %v", err)
	}

	queue, err := DefineQueue(channel, "test_queue")

	if err != nil {
		t.Fatalf("There was a problem on defining queue: %v", err)
	}

	testBody := "{\"indices\": [\"^BVSP\", \"AAPL\"]}"
	channel.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(testBody),
		},
	)

	_, err = OpenSubscriber(channel, queue.Name)

	if err != nil {
		t.Fatalf("OpenSubscriber() should return a channel for messages, but returned error: %v", err)
	}

	os.Setenv("AMQP_USERNAME", "")
	os.Setenv("AMQP_PASSWORD", "")
	os.Setenv("AMQP_DEFAULT_PORT", "")
}

func TestHandleReceivedIndices(t *testing.T) {
	os.Setenv("AMQP_USERNAME", "guest")
	os.Setenv("AMQP_PASSWORD", "guest")
	os.Setenv("AMQP_DEFAULT_PORT", "5672")

	connection, err := OpenConnection()
	defer connection.Close()

	if err != nil {
		t.Fatal(
			"There was a problem on opening connection for testing. Maybe RabbitMQ server is down ?! This is an integration test, so it's necessary to have an opened connection to AMQP server.",
		)
	}

	channel, err := OpenChannel(connection)
	defer channel.Close()

	if err != nil {
		t.Fatalf("There was a problem on opening channel: %v", err)
	}

	queue, err := DefineQueue(channel, "test_queue")

	if err != nil {
		t.Fatalf("There was a problem on defining queue: %v", err)
	}

	testBody := "{\"indices\": [\"^BVSP\", \"AAPL\"]}"
	channel.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(testBody),
		},
	)

	subscriber, err := channel.Consume(
		queue.Name, "", true, false, false, false, nil,
	)

	if err != nil {
		t.Fatal()
	}

	newArray := []string{}

	HandleReceivedIndices(
		subscriber, func(arr []string) {
			for _, item := range arr {
				newArray = append(newArray, item)
			}
		},
	)

	if strings.Join(newArray, ",") != "^BVSP,AAPL" {
		t.Fatal("Should handle subscriber received indices with handler, but nothing happened.")
	}

	os.Setenv("AMQP_USERNAME", "")
	os.Setenv("AMQP_PASSWORD", "")
	os.Setenv("AMQP_DEFAULT_PORT", "")
}

func TestPublishIndices(t *testing.T) {
	os.Setenv("AMQP_USERNAME", "guest")
	os.Setenv("AMQP_PASSWORD", "guest")
	os.Setenv("AMQP_DEFAULT_PORT", "5672")

	connection, err := OpenConnection()
	defer connection.Close()

	if err != nil {
		t.Fatal(
			"There was a problem on opening connection for testing. Maybe RabbitMQ server is down ?! This is an integration test, so it's necessary to have an opened connection to AMQP server.",
		)
	}

	channel, err := OpenChannel(connection)
	defer channel.Close()

	if err != nil {
		t.Fatalf("There was a problem on opening channel: %v", err)
	}

	queue, err := DefineQueue(channel, "test_queue")

	if err != nil {
		t.Fatalf("There was a problem on defining queue: %v", err)
	}

	result := &exchange.ExchangesResult{
		Exchanges: map[string]exchange.Exchange{
			"Nikkei 225":    exchange.Exchange{Name: "Nikkei 225", Symbol: "^n225", PercentChange: "-0.91%", ChangeInPoints: "-172.98", LastTradeDate: "4/14/2017", LastTradeTime: "3:15pm"},
			"Alphabet Inc.": exchange.Exchange{Name: "Alphabet Inc.", Symbol: "GOOGL", PercentChange: "-0.09%", ChangeInPoints: "-0.76", LastTradeDate: "4/13/2017", LastTradeTime: "4:00pm"},
		},
	}

	publishError := PublishIndices(channel, queue.Name, result)

	if publishError != nil {
		t.Fatalf("PublishIndices should run without errors, but an error was raised: %v", err)
	}

	msgs, err := channel.Consume(
		queue.Name, "", true, false, false, false, nil,
	)

	expected := "{\"Alphabet Inc.\":{\"Name\":\"Alphabet Inc.\",\"Symbol\":\"GOOGL\",\"PercentChange\":\"-0.09%\",\"ChangeInPoints\":\"-0.76\",\"LastTradeDate\":\"4/13/2017\",\"LastTradeTime\":\"4:00pm\"},\"Nikkei 225\":{\"Name\":\"Nikkei 225\",\"Symbol\":\"^n225\",\"PercentChange\":\"-0.91%\",\"ChangeInPoints\":\"-172.98\",\"LastTradeDate\":\"4/14/2017\",\"LastTradeTime\":\"3:15pm\"}}"

	for msg := range msgs {
		parsedBody := string(msg.Body)

		if parsedBody != expected {
			t.Fatalf("Published body should be %s, but it is %s", expected, parsedBody)
		}
	}

	os.Setenv("AMQP_USERNAME", "")
	os.Setenv("AMQP_PASSWORD", "")
	os.Setenv("AMQP_DEFAULT_PORT", "")
}
